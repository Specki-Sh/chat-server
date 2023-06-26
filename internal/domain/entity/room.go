package entity

type Room struct {
	ID      int       `json:"id"`
	OwnerID int       `json:"owner_id"`
	Name    string    `json:"name"`
	Members []*Member `json:"members"`
}
