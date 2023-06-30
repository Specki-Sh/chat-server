package entity

type Room struct {
	ID      ID        `json:"id"`
	OwnerID ID        `json:"owner_id"`
	Name    string    `json:"name"`
	Members []*Member `json:"members"`
}

type CreateRoomReq struct {
	OwnerID ID     `json:"owner_id"`
	Name    string `json:"name"`
}

type CreateRoomRes struct {
	ID      ID     `json:"id"`
	OwnerID ID     `json:"owner_id"`
	Name    string `json:"name"`
}

type EditRoomReq struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}

type EditRoomRes struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}
