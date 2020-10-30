package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
	clientURI := "mongodb+srv://" + os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@cluster0.kljzg.mongodb.net/foo?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(clientURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		fmt.Println("Incorrect clientURI\n Maybe env variables aren't set?")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Couldn't connect to the database: ", err)
	}
	fmt.Println("Connected to MongoDB!")
	// defer client.Disconnect(ctx)
	return client
}
