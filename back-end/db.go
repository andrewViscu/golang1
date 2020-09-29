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
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	defer client.Disconnect(ctx)
	return client
}