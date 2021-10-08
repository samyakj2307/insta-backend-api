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
	UserId   primitive.ObjectID `bson:"user_id,omitempty"`
	Name     string             `json:"name,omitempty"`
	Email    string             `json:"email,omitempty"`
	Password string             `json:"password,omitempty"`
}

type Users struct {
	AllUsers []User
}

type Post struct {
	PostId          primitive.ObjectID `json:"post_id,omitempty"`
	Caption         string             `json:"caption"`
	ImageUrl        string             `json:"image_url,omitempty"`
	PostedTimestamp time.Time          `json:"posted_timestamp,omitempty"`
	UserId          primitive.ObjectID `json:"user_id,omitempty"`
}

type Posts struct {
	AllPosts []Post
}
