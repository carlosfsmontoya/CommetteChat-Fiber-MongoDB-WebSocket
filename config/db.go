package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		log.Println("Using default values")
	}
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI no está configurado")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	DB = client
	log.Println("Connected to MongoDB!")
}

func DisconnectDB() {
	if err := DB.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Error desconectando de MongoDB: %v", err)
	}
	log.Println("Desconectado de MongoDB!")
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("MONGODB_DBNAME")
	if dbName == "" {
		log.Fatal("MONGODB_DBNAME no está configurado")
	}
	return client.Database(dbName).Collection(collectionName)
}
