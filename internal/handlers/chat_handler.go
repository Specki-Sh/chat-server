package handlers

import (
	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	"chat-server/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"nhooyr.io/websocket"
	"strconv"
	"sync"
)

type ChatHandler struct {
	messageUseCase use_case.MessageUseCase

	chatsMu sync.Mutex
	chats   map[int]*service.Chat

	messageBuffSize   int
	broadcastBuffSize int
}

func NewChatHandler(messageUseCase use_case.MessageUseCase) *ChatHandler {
	return &ChatHandler{
		messageUseCase: messageUseCase,
		chats:          make(map[int]*service.Chat),
	}
}

func (ch *ChatHandler) JoinRoom(c *gin.Context) {
	// Get the chat ID from the "id" URL parameter.
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(c.Query("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Accept the WebSocket connection.
	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "")

	// Create a new client and add it to the chat.
	cl := service.NewClient(conn, ch.messageBuffSize, roomID, userID)
	ch.addClient(roomID, cl)
	defer ch.deleteClient(roomID, cl)

	// Start a go-routine to send messages to the client.
	go cl.WriteMessage()

	// Start the broadcast manager for the chat (if it is not already started).
	go ch.startBroadcastManager(roomID)

	// Read messages from the client and send them to the chat's broadcast channel.
	cl.ReadMessage(ch.chats[roomID].Broadcast)
}

func (ch *ChatHandler) EditMessage(c *gin.Context) {
	var req use_case.EditMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	req.ID = id
	message, err := ch.messageUseCase.EditMessageContent(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, message)
}

func (ch *ChatHandler) DeleteMessage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	err = ch.messageUseCase.RemoveMessageByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (ch *ChatHandler) DeleteAllMessageFromRoom(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	err = ch.messageUseCase.RemoveMessagesByRoomID(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (ch *ChatHandler) GetMessagesPaginate(c *gin.Context) {
	var req use_case.GetMessagesPaginateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomID, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	req.RoomID = roomID
	messages, err := ch.messageUseCase.GetMessagesPaginate(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (ch *ChatHandler) MessagePermissionMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		messageID, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
			return
		}

		userID, err := getUserID(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		isOwner, err := ch.messageUseCase.IsMessageOwner(userID, messageID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !isOwner {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}

func (ch *ChatHandler) BroadcastMessageUpdateMiddleware(c *gin.Context) {
	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
		return
	}

	c.Next()
	if c.IsAborted() {
		return
	}

	msg, err := ch.messageUseCase.GetMessageByID(messageID)
	if err != nil {
		return
	}
	ch.sendMessageForAllClientInRoom(msg)
}

// startBroadcastManager starts the broadcast manager for the chat with the specified ID (if it is not already started).
func (ch *ChatHandler) startBroadcastManager(roomID int) {
	if chat, ok := ch.chats[roomID]; ok && !chat.BroadcastManagerStatus {
		ch.broadcastManager(chat.Broadcast)
	}
}

// broadcastManager handles messages from the broadcast channel and sends them to all clients in the chat.
func (ch *ChatHandler) broadcastManager(broadcast chan *entity.Message) {
	for {
		select {
		case msg := <-broadcast:
			req := use_case.NewCreateMessageReq(msg)
			message, err := ch.messageUseCase.CreateMessage(req)
			if err != nil {
				// log
				continue
			}
			ch.sendMessageForAllClientInRoom(message)
		}
	}
}

func (ch *ChatHandler) sendMessageForAllClientInRoom(msg *entity.Message) {
	if chat, ok := ch.chats[msg.RoomID]; ok {
		for cl := range chat.Clients {
			cl.Message <- msg
		}
	}
}

// addClient adds a client to the chat with the specified ID. If the chat does not exist, it is created.
func (ch *ChatHandler) addClient(roomID int, c *service.Client) {
	ch.chatsMu.Lock()
	defer ch.chatsMu.Unlock()

	chat, ok := ch.chats[roomID]
	if !ok {
		chat = service.NewChat(ch.broadcastBuffSize)
		ch.chats[roomID] = chat
	}

	chat.Clients[c] = struct{}{}
}

// deleteClient removes a client from the chat with the specified ID. If there are no more clients in the chat, then the chat is deleted.
func (ch *ChatHandler) deleteClient(roomID int, c *service.Client) {
	ch.chatsMu.Lock()
	delete(ch.chats[roomID].Clients, c)
	if len(ch.chats[roomID].Clients) == 0 {
		delete(ch.chats, roomID)
	}
	ch.chatsMu.Unlock()
}
