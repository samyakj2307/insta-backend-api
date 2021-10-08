package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

import . "AppointyTask/models"

var dbHandler = BaseHandler{}

//var allUsers = Users{}

//var allPosts = Posts{}

var ctx context.Context

func makeDBConnection() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
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
	}(client, ctx)

	dbHandler.Db = client.Database("appointyInstaDB")

	dbHandler.Users = dbHandler.Db.Collection("users")
	dbHandler.Posts = dbHandler.Db.Collection("posts")

	http.HandleFunc("/users", addUser)
	http.HandleFunc("/users/", getUser)

	http.ListenAndServe(":8080", nil)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Add("Content-Type", "application/json")
		userId, _ := primitive.ObjectIDFromHex(strings.Split(r.URL.String(), "/users/")[1])

		cursor, err := dbHandler.Users.Find(ctx, bson.M{"_id": userId})

		if err != nil {
			log.Fatal(err)
		}
		var currentUser []bson.M
		if err = cursor.All(ctx, &currentUser); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(currentUser[0])

	} else {
		w.Write([]byte("Method Not Allowed: " + r.Method))
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)

		var user User

		err := decoder.Decode(&user)
		if err != nil {
			panic(err)
		}

		hashedPassword, passwordHashingErr := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

		if passwordHashingErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user.Password = string(hashedPassword)

		one, err := dbHandler.Users.InsertOne(ctx, user)
		if err != nil {
			return
		}

		json.NewEncoder(w).Encode(one)
	} else {
		w.Write([]byte("Method Not Allowed: " + r.Method))
	}
}
