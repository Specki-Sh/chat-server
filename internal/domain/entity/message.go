package entity

import "time"

type Message struct {
	ID        ID             `json:"id"`
	SenderID  ID             `json:"sender_id"`
	RoomID    ID             `json:"room_id"`
	Content   NonEmptyString `json:"content"`
	Status    string         `json:"status"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at"`
	IsActive  bool           `json:"is_active"`
}

type CreateMessageReq struct {
	SenderID ID             `json:"sender_id"`
	RoomID   ID             `json:"room_id"`
	Content  NonEmptyString `json:"content"`
}

func NewCreateMessageReq(message *Message) *CreateMessageReq {
	return &CreateMessageReq{
		SenderID: message.SenderID,
		RoomID:   message.RoomID,
		Content:  message.Content,
	}
}

type EditMessageReq struct {
	ID      ID             `json:"id"`
	Content NonEmptyString `json:"content"`
}

func (e *EditMessageReq) Validate() error {
	if err := e.ID.Validate(); err != nil {
		return err
	}
	if err := e.Content.Validate(); err != nil {
		return err
	}
	return nil
}

type GetMessagesPaginateReq struct {
	RoomID  ID   `json:"room_id"`
	PerPage uint `json:"per_page"`
	Page    uint `json:"page"`
}

func (g *GetMessagesPaginateReq) Validate() error {
	if err := g.RoomID.Validate(); err != nil {
		return err
	}
	return nil
}
