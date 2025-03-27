package main

import (
	"fmt"
	"log"
	"noteapp-framework-backend/config"
	"noteapp-framework-backend/handlers"
)

func main() {
	// Initialize the database
	config.DBInit()

	// Create a new note
	err := handlers.CreateNote("My First Note", "This is the content of the note", 1)
	if err != nil {
		log.Fatalf("Failed to create note: %v", err)
	}
	fmt.Println("Note created successfully!")

	// Fetch all notes
	notes, err := handlers.GetNotes()
	if err != nil {
		log.Fatalf("Failed to fetch notes: %v", err)
	}
	fmt.Printf("Notes: %+v\n", notes)
}
