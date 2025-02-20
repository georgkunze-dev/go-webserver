package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB repräsentiert die Verbindung zu der MongoDB-Datenbank
type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

// Connect verbindet sich mit MongoDB und gibt eine MongoDB-Instanz zurück.
func Connect() (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Konfiguration der MongoDB-Verbindung
	opt := options.Client().ApplyURI("mongodb://mongodb:27017")
	// Verbindung herstellen
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}

	// Verbindung überprüfen
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	database := client.Database("mydb")
	log.Println(">>> Verbindung zur Datenbank hergestellt")
	return &MongoDB{client: client, db: database}, nil
}

// Close trennt die Verbindung zur MongoDB-Datenbank.
func (db *MongoDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db.client.Disconnect(ctx)
}

// SaveFile speichert eine Datei in der MongoDB-Datenbank.
func (db *MongoDB) SaveFile(filename string, data []byte) error {
	// Initialisierung der DB prüfen
	if db.db == nil {
		log.Println("Datenbank nicht initialisiert")
	}
	collection := db.db.Collection("files")

	// Filter erstellen
	filter := bson.D{{"filename", filename}}
	// BSON-Dokument erstellen
	entry := bson.D{{"$set", bson.D{{"filename", filename}, {"data", data}}}}

	// Datei in die MongoDB-Datenbank speichern
	opt := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, entry, opt)
	if err != nil {
		return fmt.Errorf("SaveFile.UpdateOne filename %s: %w", filename, err)
	}
	return nil
}

// GetFile ruft eine Datei basierend auf dem Dateinamen aus der MongoDB-Datenbank ab.
func (db *MongoDB) GetFile(filename string) ([]byte, error) {
	collection := db.db.Collection("files")
	var result struct {
		Data []byte `bson:"data"`
	}
	// Datei aus der MongoDB-Datenbank holen
	err := collection.FindOne(context.Background(), bson.M{"filename": filename}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("GetFile.FindOne: %w", err)
	}
	return result.Data, nil
}
