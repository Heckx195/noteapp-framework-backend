package models

type Note struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	NotebookID int    `json:"notebook_id"`
}
