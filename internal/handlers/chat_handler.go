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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "")

	cl := service.NewClient(conn, h.messageBuffSize, id)
	h.addClient(id, cl)
	defer h.deleteClient(id, cl)

	go cl.WriteMessage()

	go h.startBroadcastManager(id)
	cl.ReadMessage(h.chats[id].Broadcast)
}

func (h *ChatHandler) startBroadcastManager(roomID int) {
	if chat, ok := h.chats[roomID]; ok && !chat.BroadcastManagerStatus {
		h.broadcastManager(chat.Broadcast)
	}
}

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

func (h *ChatHandler) deleteClient(roomId int, c *service.Client) {
	h.chatsMu.Lock()
	delete(h.chats[roomId].Clients, c)
	if len(h.chats[roomId].Clients) == 0 {
		delete(h.chats, roomId)
	}
	h.chatsMu.Unlock()
}
