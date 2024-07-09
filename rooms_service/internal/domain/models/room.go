package models

type Room struct {
	ID       uint64 `json:"id"`
	Name     string `json:"room_name"`
	Capacity uint64 `json:"capacity"`
}
