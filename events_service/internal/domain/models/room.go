package models

type Room struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	Capacity uint64 `json:"capacity"`
}
