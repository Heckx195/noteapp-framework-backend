package handlers

import (
    "noteapp-framework-backend/config"
    "noteapp-framework-backend/models"
)

// CreateNote inserts a new note into the database
func CreateNote(title, content string, notebookID int) error {
    query := "INSERT INTO notes (title, content, notebook_id) VALUES ($1, $2, $3)"
    _, err := config.DB.Exec(query, title, content, notebookID)
    return err
}

func GetNotes() ([]models.Note, error) {
    query := "SELECT id, title, content, notebook_id FROM notes"
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notes []models.Note
    for rows.Next() {
        var note models.Note
        err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.NotebookID)
        if err != nil {
            return nil, err
        }
        notes = append(notes, note)
    }
    return notes, nil
}