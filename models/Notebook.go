package models

type Notebook struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID uint   `json:"user_id"`
}
