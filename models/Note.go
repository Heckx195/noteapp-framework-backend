package models

type Note struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	NotebookID uint   `json:"notebook_id"`
	UserID     uint   `json:"user_id"`
}
