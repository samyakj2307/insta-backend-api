package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

import . "AppointyTask/models"
import . "AppointyTask/api"

var dbHandler = BaseHandler{}

func makeDBConnection() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func main() {

	client := makeDBConnection()

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Println("Error Disconnecting Database")
		}
	}(client, context.Background())

	dbHandler.Db = client.Database("appointyInstaDB")

	dbHandler.Users = dbHandler.Db.Collection("users")
	dbHandler.Posts = dbHandler.Db.Collection("posts")
	http.HandleFunc("/users", addUserEndpoint)
	http.HandleFunc("/users/", getUserEndpoint)

	http.HandleFunc("/posts", addPostEndpoint)
	http.HandleFunc("/posts/", getPostEndpoint)
	http.HandleFunc("/posts/users/", getAllPostsOfUserEndpoint)

	http.ListenAndServe(":8080", nil)
}

func addUserEndpoint(w http.ResponseWriter, r *http.Request) {
	AddUser(w, r, dbHandler)
}

func getUserEndpoint(w http.ResponseWriter, r *http.Request) {
	GetUser(w, r, dbHandler)
}

func addPostEndpoint(w http.ResponseWriter, r *http.Request) {
	AddPost(w, r, dbHandler)
}

func getPostEndpoint(w http.ResponseWriter, r *http.Request) {
	GetPost(w, r, dbHandler)
}

func getAllPostsOfUserEndpoint(w http.ResponseWriter, r *http.Request) {
	GetAllPostsOfUser(w, r, dbHandler)
}
