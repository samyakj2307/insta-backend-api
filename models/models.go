package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type BaseHandler struct {
	Db    *mongo.Database
	Ctx   *context.Context
	Users *mongo.Collection
	Posts *mongo.Collection
}

type User struct {
	_id      primitive.ObjectID `bson:"user_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
}

type Users struct {
	AllUsers []User
}

type Post struct {
	_id             primitive.ObjectID `bson:"post_id,omitempty"`
	Caption         string             `bson:"caption"`
	ImageUrl        string             `bson:"image_url,omitempty"`
	PostedTimestamp time.Time          `bson:"posted_timestamp,omitempty"`
	UserId          primitive.ObjectID `bson:"user_id,omitempty"`
}

type Posts struct {
	AllPosts []Post
}
