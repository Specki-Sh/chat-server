package entity

import "time"

type Message struct {
	ID        int        `json:"id"`
	SenderID  int        `json:"sender_id"`
	RoomID    int        `json:"room_id"`
	Content   string     `json:"content"`
	Status    string     `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	IsActive  bool       `json:"is_active"`
}
