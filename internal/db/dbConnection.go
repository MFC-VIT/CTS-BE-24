package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var DB *mongo.Database

func ConnectMongoDB(uri string,dbName string) (*mongo.Client, *mongo.Database){
	clientOptions:= options.Client().ApplyURI(uri)
	client,err:= mongo.Connect(context.Background(),clientOptions)
	if err !=nil{
		log.Fatalf("Failed to connect to MongoDB: %v",err)
	}
	err = client.Ping(context.Background(),nil)
	if err!=nil{
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db := client.Database(dbName)
	log.Println("Connected to MongoDB!")
	return client,db
}