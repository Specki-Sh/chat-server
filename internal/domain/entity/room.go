package entity

type Room struct {
	ID      int       `json:"id"`
	OwnerID int       `json:"owner_id"`
	Name    string    `json:"name"`
	Members []*Member `json:"members"`
}

type CreateRoomReq struct {
	OwnerID int    `json:"owner_id"`
	Name    string `json:"name"`
}

type CreateRoomRes struct {
	ID      int    `json:"id"`
	OwnerID int    `json:"owner_id"`
	Name    string `json:"name"`
}

type EditRoomReq struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type EditRoomRes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
