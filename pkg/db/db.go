package db

import (
	"os"
	"fmt"
	"log"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBConnect() *mongo.Client {
	clientURI := "mongodb+srv://" + os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@cluster0.kljzg.mongodb.net/foo?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(clientURI)
	client, err := mongo.NewClient(clientOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Couldn't connect to the database: ",err)
	}
	fmt.Println("Connected to MongoDB!")
	// defer client.Disconnect(ctx)
	return client
}