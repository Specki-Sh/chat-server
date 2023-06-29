package entity

import "time"

type Message struct {
	ID        ID         `json:"id"`
	SenderID  ID         `json:"sender_id"`
	RoomID    ID         `json:"room_id"`
	Content   string     `json:"content"`
	Status    string     `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	IsActive  bool       `json:"is_active"`
}

type CreateMessageReq struct {
	SenderID ID     `json:"sender_id"`
	RoomID   ID     `json:"room_id"`
	Content  string `json:"content"`
}

func NewCreateMessageReq(message *Message) *CreateMessageReq {
	return &CreateMessageReq{
		SenderID: message.SenderID,
		RoomID:   message.RoomID,
		Content:  message.Content,
	}
}

type EditMessageReq struct {
	ID      ID     `json:"id"`
	Content string `json:"content"`
}

type GetMessagesPaginateReq struct {
	RoomID  ID  `json:"room_id"`
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
}
