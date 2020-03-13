package backend

import (
	"context"
	"fmt"
	"log"

	"github.com/gargath/mongoose/pkg/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) ListUsers() ([]*entities.User, error) {
	users := []*entities.User{}
	coll := b.m.Database("test").Collection("Users")

	cur, err := coll.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		return users, fmt.Errorf("failed to find users: %v", err)
	}

	for cur.Next(context.TODO()) {
		var u entities.User
		err := cur.Decode(&u)
		if err != nil {
			return users, fmt.Errorf("failed to decode Mongo document for user: %v", err)
		}
		users = append(users, &u)
	}

	if err := cur.Err(); err != nil {
		return users, fmt.Errorf("encountered error while reading documents from Mongo: %v", err)
	}

	cur.Close(context.TODO())

	return users, nil
}

func (b *Backend) InsertUser(u *entities.User) (string, error) {
	coll := b.m.Database("test").Collection("Users")
	insertResult, err := coll.InsertOne(context.TODO(), u)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %v", err)
	}
	log.Printf("inserted new user: %+v", insertResult)
	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}
