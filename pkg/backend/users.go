package backend

import (
	"context"
	"fmt"

	"github.com/gargath/mongoose/pkg/objects"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) ListUsers() ([]*objects.User, error) {
	users := []*objects.User{}
	coll := b.m.Database("test").Collection("Users")

	cur, err := coll.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		return users, fmt.Errorf("failed to find users: %v", err)
	}

	for cur.Next(context.TODO()) {
		var u objects.User
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
