package entity

import "time"

type Message struct {
	ID        int       `json:"id"`
	SenderID  int       `json:"sender_id"`
	RoomID    int       `json:"room_id"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
