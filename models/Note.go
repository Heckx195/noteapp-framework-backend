package models

import (
	"time"
)

type Note struct {
	ID         string    `bson:"_id,omitempty"`
	Title      string    `bson:"title"`
	Content    string    `bson:"content"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
	NotebookID string    `bson:"notebook_id"`
}
