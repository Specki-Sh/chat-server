package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"nhooyr.io/websocket"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	"chat-server/internal/service"
)

type ChatHandler struct {
	messageUseCase use_case.MessageUseCase

	chatsMu sync.Mutex
	chats   map[entity.ID]*service.Chat

	messageBuffSize   int
	broadcastBuffSize int
}

func NewChatHandler(messageUseCase use_case.MessageUseCase, logger *logrus.Logger) *ChatHandler {
	return &ChatHandler{
		messageUseCase: messageUseCase,
		chats:          make(map[entity.ID]*service.Chat),
	}
}

func (ch *ChatHandler) JoinRoom(c *gin.Context) {
	roomID, userID, err := ch.getRoomIDAndUserIDParams(c)
	if err != nil {
		log.Printf("error getting params: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := ch.acceptWebSocket(c)
	if err != nil {
		log.Printf("error accepting WebSocket connection: %v", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "")

	cl := ch.createClient(conn, roomID, userID)
	defer ch.deleteClient(roomID, cl)

	go cl.WriteMessage()
	go ch.startBroadcastManager(roomID)

	log.Printf("user joined room: %d %d", userID, roomID)

	cl.ReadMessage(ch.chats[roomID].Broadcast)
}

func (ch *ChatHandler) getRoomIDAndUserIDParams(c *gin.Context) (entity.ID, entity.ID, error) {
	roomIDInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, 0, fmt.Errorf("error converting roomID to int: %w", err)
	}

	userIDInt, err := getUserID(c)
	if err != nil {
		return 0, 0, fmt.Errorf("error converting userID to int: %w", err)
	}

	return entity.ID(roomIDInt), entity.ID(userIDInt), nil
}

func (ch *ChatHandler) acceptWebSocket(c *gin.Context) (*websocket.Conn, error) {
	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("error accepting WebSocket connection: %w", err)
	}
	return conn, nil
}

func (ch *ChatHandler) createClient(
	conn *websocket.Conn,
	roomID entity.ID,
	userID entity.ID,
) *service.Client {
	cl := service.NewClient(conn, ch.messageBuffSize, roomID, userID)
	ch.addClient(roomID, cl)
	return cl
}

func (ch *ChatHandler) EditMessage(c *gin.Context) {
	var req entity.EditMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("error binding JSON: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("error converting message ID to int: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	req.ID = entity.ID(id)

	if err := req.Validate(); err != nil {
		log.Printf("error validating request: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := ch.messageUseCase.EditMessageContent(&req)
	if err != nil {
		log.Printf("error editing message content: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("message edited: %d", id)
	c.JSON(http.StatusOK, message)
}

func (ch *ChatHandler) DeleteMessage(c *gin.Context) {
	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("error converting message ID to int: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	id := entity.ID(idInt)

	err = ch.messageUseCase.RemoveMessageByID(id)
	if err != nil {
		log.Printf("error removing message by ID: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("message deleted: %d", id)
	c.Status(http.StatusNoContent)
}

func (ch *ChatHandler) DeleteAllMessageFromRoom(c *gin.Context) {
	roomIDInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("error converting room ID to int: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	roomID := entity.ID(roomIDInt)

	err = ch.messageUseCase.RemoveMessageBulkByRoomID(roomID)
	if err != nil {
		log.Printf("error removing messages by room ID: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("all messages deleted from room: %d", roomID)
	c.Status(http.StatusNoContent)
}

func (ch *ChatHandler) GetMessageBulkPaginate(c *gin.Context) {
	var req entity.GetMessageBulkPaginateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("error binding JSON: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomID, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		log.Printf("error converting room ID to int: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}
	req.RoomID = entity.ID(roomID)

	if err := req.Validate(); err != nil {
		log.Printf("error validating request: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messages, err := ch.messageUseCase.GetMessageBulkPaginate(&req)
	if err != nil {
		log.Printf("error getting messages paginate: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("messages paginate retrieved: %d", roomID)
	c.JSON(http.StatusOK, messages)
}

func (ch *ChatHandler) MessagePermissionMiddlewareByParam(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		messageIDInt, err := strconv.Atoi(c.Param(paramKey))
		if err != nil {
			log.Printf("error converting message ID to int: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
			return
		}
		messageID := entity.ID(messageIDInt)

		userID, err := getUserID(c)
		if err != nil {
			log.Printf("error getting user ID: %v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		isOwner, err := ch.messageUseCase.IsMessageOwner(userID, messageID)
		if err != nil {
			log.Printf("error checking if user is message owner: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !isOwner {
			log.Printf("access denied to message: %d %d", userID, messageID)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		log.Printf("access granted to message: %d %d", userID, messageID)
		c.Next()
	}
}

func (ch *ChatHandler) BroadcastMessageUpdateMiddleware(c *gin.Context) {
	messageIDInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("error converting message ID to int: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
		return
	}
	messageID := entity.ID(messageIDInt)

	c.Next()
	if c.IsAborted() {
		return
	}

	msg, err := ch.messageUseCase.GetMessageByID(messageID)
	if err != nil {
		log.Printf("error getting message by ID: %v", err)
		return
	}
	ch.sendMessageForAllClientInRoom(msg)
}

// startBroadcastManager starts the broadcast manager for the chat with the specified ID (if it is not already started).
func (ch *ChatHandler) startBroadcastManager(roomID entity.ID) {
	if chat, ok := ch.chats[roomID]; ok && !chat.BroadcastManagerStatus {
		ch.broadcastManager(chat.Broadcast)
	}
}

// broadcastManager handles messages from the broadcast channel and sends them to all clients in the chat.
func (ch *ChatHandler) broadcastManager(broadcast chan *entity.Message) {
	for {
		select {
		case msg := <-broadcast:
			req := entity.NewCreateMessageReq(msg)
			message, err := ch.messageUseCase.CreateMessage(req)
			if err != nil {
				log.Printf("error creating message: %v", err)
				continue
			}
			ch.sendMessageForAllClientInRoom(message)
			log.Printf("message broadcasted: %d", message.ID)
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
func (ch *ChatHandler) addClient(roomID entity.ID, c *service.Client) {
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
func (ch *ChatHandler) deleteClient(roomID entity.ID, c *service.Client) {
	ch.chatsMu.Lock()
	delete(ch.chats[roomID].Clients, c)
	if len(ch.chats[roomID].Clients) == 0 {
		delete(ch.chats, roomID)
	}
	ch.chatsMu.Unlock()
}
