package service

import (
	"chat-server/internal/domain/entity"
	"context"
	"encoding/json"
	"io"
	"nhooyr.io/websocket"
)

type Chat struct {
	Broadcast              chan *entity.Message
	BroadcastManagerStatus bool
	Clients                map[*Client]struct{}
}

func NewChat(bcBuffSize int) *Chat {
	return &Chat{
		Broadcast: make(chan *entity.Message, bcBuffSize),
		Clients:   make(map[*Client]struct{}),

		BroadcastManagerStatus: false,
	}
}

type Client struct {
	Conn    *websocket.Conn
	Message chan *entity.Message
	RoomID  entity.ID
	UserID  entity.ID
}

func NewClient(conn *websocket.Conn, messageBuffSize int, roomID entity.ID, userID entity.ID) *Client {
	return &Client{
		Conn:    conn,
		Message: make(chan *entity.Message, messageBuffSize),
		RoomID:  roomID,
		UserID:  userID,
	}
}

func (c *Client) WriteMessage() {
	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		messageBytes, err := json.Marshal(message)
		if err != nil {
			return
		}

		if err := c.Conn.Write(context.Background(), websocket.MessageText, messageBytes); err != nil {
			return
		}
	}
}

func (c *Client) ReadMessage(broadcast chan *entity.Message) {
	for {
		_, m, err := c.Conn.Read(context.Background())
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				err == io.EOF {
				return
			}
			return
		}

		msg := &entity.Message{
			RoomID:   c.RoomID,
			SenderID: c.UserID,
			Content:  string(m),
		}

		broadcast <- msg
	}
}
