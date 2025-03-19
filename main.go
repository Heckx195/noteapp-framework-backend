package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"noteapp-framework-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for your collection
	collection := client.Database("testdb").Collection("notes")

	// Create a new note
	note := models.Note{
		Title:      "My Note",
		Content:    "This is a note",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		NotebookID: "notebook1",
	}

	// Insert the note into the collection
	insertResult, err := collection.InsertOne(context.TODO(), note)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// Find a single note
	var result models.Note
	filter := bson.D{{Key: "title", Value: "My Note"}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)
}
