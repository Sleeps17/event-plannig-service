package models

type Room struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	Capacity uint32 `json:"capacity"`
}
