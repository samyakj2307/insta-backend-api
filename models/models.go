package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/url"
	"time"
)

type BaseHandler struct {
	Db    *mongo.Database
	Ctx   *context.Context
	Users *mongo.Collection
	Posts *mongo.Collection
}

type User struct {
	_id      primitive.ObjectID `bson:"user_id"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type Users struct {
	AllUsers []User
}

type Post struct {
	_id             primitive.ObjectID `bson:"post_id"`
	Caption         string             `bson:"caption"`
	ImageUrl        url.URL            `bson:"image_url"`
	PostedTimestamp time.Time          `bson:"posted_timestamp"`
	UserId          primitive.ObjectID `bson:"user_id"`
}

type Posts struct {
	AllPosts []Post
}
