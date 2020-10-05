package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connectionURI = "mongodb+srv://andrewViscu:AquaDash1324@cluster0.kljzg.mongodb.net/foo?retryWrites=true&w=majority"

func db() *mongo.Client {
	clientOptions := options.Client().ApplyURI(connectionURI)
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