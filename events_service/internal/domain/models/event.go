package models

import "time"

type Event struct {
	ID           uint64      `json:"id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	StartDate    time.Time   `json:"start_date"`
	EndDate      time.Time   `json:"end_date"`
	Room         *Room       `json:"room"`
	Creator      *Employee   `json:"creator"`
	Participants []*Employee `json:"participants"`
}
