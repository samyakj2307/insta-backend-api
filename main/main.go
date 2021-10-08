package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

import . "AppointyTask/models"

var dbInstance *mongo.Database
var userCollection *mongo.Collection
var postCollection *mongo.Collection
var contextInstance *context.Context
var dbHandler = BaseHandler{Db: dbInstance, Users: userCollection, Posts: postCollection, Ctx: contextInstance}

var users []User
var allUsers = Users{AllUsers: users}

var posts []Post
var allPosts = Posts{AllPosts: posts}

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

	//http.HandleFunc("/getUsers", getUser)
	http.HandleFunc("/addUsers", addUser)

	http.ListenAndServe(":8080", nil)
}

func addUser(w http.ResponseWriter, r *http.Request) {

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

	log.Println(one)
	resErr := json.NewEncoder(w).Encode(one)
	if resErr != nil {
		return
	}
}
