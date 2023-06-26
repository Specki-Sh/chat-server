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

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		chats: make(map[int]*service.Chat),
	}
}

func (h *ChatHandler) JoinRoom(c *gin.Context) {
	// Get the chat ID from the "id" URL parameter.
	id, err := strconv.Atoi(c.Param("id"))
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
	cl := service.NewClient(conn, h.messageBuffSize, id)
	h.addClient(id, cl)
	defer h.deleteClient(id, cl)

	// Start a go-routine to send messages to the client.
	go cl.WriteMessage()

	// Start the broadcast manager for the chat (if it is not already started).
	go h.startBroadcastManager(id)

	// Read messages from the client and send them to the chat's broadcast channel.
	cl.ReadMessage(h.chats[id].Broadcast)
}

func (h *ChatHandler) EditMessage(c *gin.Context) {
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
	message, err := h.messageUseCase.EditMessageContent(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, message)
}

func (h *ChatHandler) DeleteMessage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	err = h.messageUseCase.RemoveMessageByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ChatHandler) DeleteAllMessageFromRoom(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	err = h.messageUseCase.RemoveMessagesByRoomID(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ChatHandler) GetMessagesPaginate(c *gin.Context) {
	var req use_case.GetMessagesPaginateReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomID, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	req.RoomID = roomID

	messages, err := h.messageUseCase.GetMessagesPaginate(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (h *ChatHandler) MessagePermissionMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		messageID, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
			return
		}

		userID, err := getUserId(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		isOwner, err := h.messageUseCase.IsMessageOwner(userID, messageID)
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

// startBroadcastManager starts the broadcast manager for the chat with the specified ID (if it is not already started).
func (h *ChatHandler) startBroadcastManager(roomID int) {
	if chat, ok := h.chats[roomID]; ok && !chat.BroadcastManagerStatus {
		h.broadcastManager(chat.Broadcast)
	}
}

// broadcastManager handles messages from the broadcast channel and sends them to all clients in the chat.
func (h *ChatHandler) broadcastManager(broadcast chan *entity.Message) {
	for {
		select {
		case msg := <-broadcast:
			req := use_case.NewCreateMessageReq(msg)
			message, err := h.messageUseCase.CreateMessage(req)
			if err != nil {
				// log
				continue
			}
			if chat, ok := h.chats[message.RoomID]; ok {
				for cl := range chat.Clients {
					cl.Message <- message
				}
			}
		}
	}
}

// addClient adds a client to the chat with the specified ID. If the chat does not exist, it is created.
func (h *ChatHandler) addClient(roomID int, c *service.Client) {
	h.chatsMu.Lock()
	defer h.chatsMu.Unlock()

	chat, ok := h.chats[roomID]
	if !ok {
		chat = service.NewChat(h.broadcastBuffSize)
		h.chats[roomID] = chat
	}

	chat.Clients[c] = struct{}{}
}

// deleteClient removes a client from the chat with the specified ID. If there are no more clients in the chat, then the chat is deleted.
func (h *ChatHandler) deleteClient(roomId int, c *service.Client) {
	h.chatsMu.Lock()
	delete(h.chats[roomId].Clients, c)
	if len(h.chats[roomId].Clients) == 0 {
		delete(h.chats, roomId)
	}
	h.chatsMu.Unlock()
}
