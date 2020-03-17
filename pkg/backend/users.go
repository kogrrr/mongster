package backend

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gargath/mongoose/pkg/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) ListUsers() ([]*entities.User, error) {
	//TODO: MUST(!) filter token before returning
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

func (b *Backend) FindUserBySub(sub string) (*entities.User, error) {
	var user *entities.User
	coll := b.m.Database("test").Collection("Users")

	filter := bson.M{"sub": sub}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	u := coll.FindOne(ctx, filter)
	if u.Err() != nil {
		if u.Err() == mongo.ErrNoDocuments {
			return user, nil
		} else {
			return user, fmt.Errorf("failed to fetch user with sub %v: %v", sub, u.Err())
		}
	}
	u.Decode(&user)
	return user, nil
}

func (b *Backend) UpdateUserToken(sub string, t *entities.Token) error {
	coll := b.m.Database("test").Collection("Users")

	opts := options.Update().SetUpsert(false)
	filter := bson.M{"sub": sub}
	update := bson.D{{"$set", bson.D{{"token", t}}}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("token update did not match existing user")
	}
	return nil
}
