package models

import (
	"time"
)

type Notebook struct {
	ID        uint      `bson:"_id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Notes     []Note    `json:"notes"`
}
