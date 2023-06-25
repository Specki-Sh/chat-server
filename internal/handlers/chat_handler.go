package handlers

import (
	"chat-server/internal/domain/entity"
	"chat-server/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"nhooyr.io/websocket"
	"strconv"
	"sync"
)

type ChatHandler struct {
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
			if chat, ok := h.chats[msg.RoomID]; ok {
				for cl := range chat.Clients {
					cl.Message <- msg
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
