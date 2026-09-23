package models

import "time"

type Job struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Payload   []byte    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
