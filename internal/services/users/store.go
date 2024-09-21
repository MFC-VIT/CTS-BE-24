package users

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
    collection *mongo.Collection
    client     *mongo.Client
    db         *mongo.Database 
}

func NewUserStore(db *mongo.Database) *Store {
    collectionName := os.Getenv("MONGO_USER_COLLECTION")
    collection := db.Collection(collectionName)
    return &Store{collection: collection,db:db}
}
