package service

import (
	"chat-server/internal/domain/entity"
	"context"
	"encoding/json"
	"io"
	"nhooyr.io/websocket"
)

type Chat struct {
	Clients map[*Client]struct{}
}

type Client struct {
	Conn    *websocket.Conn
	Message chan *entity.Message
	RoomID  int
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
			RoomID:  c.RoomID,
			Content: string(m),
		}

		broadcast <- msg
	}
}
