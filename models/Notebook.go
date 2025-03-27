package models

type Notebook struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Notes []Note `json:"notes"`
}
