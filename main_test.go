package main

import (
	"bytes"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAddAndGetUser(t *testing.T) {

	//Making DB Connection
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

	//Post Data in User Collection
	postUserDetails := map[string]interface{}{
		"name":     "Samyak",
		"email":    "jainsamyak@gmail.com",
		"password": "12345678",
	}

	body, _ := json.Marshal(postUserDetails)
	postUserReq := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))

	postUserWriter := httptest.NewRecorder()
	addUserEndpoint(postUserWriter, postUserReq)
	postUserResponse := postUserWriter.Result()

	defer postUserResponse.Body.Close()

	postUserData, err := ioutil.ReadAll(postUserResponse.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var result map[string]interface{}

	json.Unmarshal([]byte(postUserData), &result)

	userIdCreated := result["InsertedID"].(string)

	// Testing if the User ID is matching with any of
	//the documents present in db or not
	// By Using Get User Method

	getUserRequest := httptest.NewRequest(http.MethodGet, "/users/"+userIdCreated, nil)

	getUserWriter := httptest.NewRecorder()
	getUserEndpoint(getUserWriter, getUserRequest)
	getUserResponse := getUserWriter.Result()

	defer getUserResponse.Body.Close()

	getUserdata, getUsererr := ioutil.ReadAll(getUserResponse.Body)
	if getUsererr != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var getUserResult map[string]interface{}

	json.Unmarshal([]byte(getUserdata), &getUserResult)

	//getUserResult["email"] = "Hello@gmail.com"

	time.Sleep(10)

	if getUserResult["email"] != postUserDetails["email"] || getUserResult["name"] != postUserDetails["name"] || getUserResult["_id"] != userIdCreated {
		t.Error("Wrong Output")
	}
}

func TestAddAndGetPost(t *testing.T) {

	//Making DB Connection
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

	//Post Data in User Collection
	addPostDetails := map[string]interface{}{
		"user_id":   "6161b07a1e6a823e09e31941",
		"image_url": "instabucket.aws.com",
		"caption":   "Hello Go Lang",
	}

	body, _ := json.Marshal(addPostDetails)
	addPostReq := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(body))

	addPostWriter := httptest.NewRecorder()
	addPostEndpoint(addPostWriter, addPostReq)
	addPostResponse := addPostWriter.Result()

	defer addPostResponse.Body.Close()

	addPostData, err := ioutil.ReadAll(addPostResponse.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var result map[string]interface{}

	json.Unmarshal([]byte(addPostData), &result)

	postIdCreated := result["InsertedID"].(string)

	// Testing if the Post ID is matching with any of
	//the documents present in db or not
	// By Using Get Post Method

	getPostRequest := httptest.NewRequest(http.MethodGet, "/posts/"+postIdCreated, nil)

	getPostWriter := httptest.NewRecorder()
	getPostEndpoint(getPostWriter, getPostRequest)
	getPostResponse := getPostWriter.Result()

	defer getPostResponse.Body.Close()

	getPostdata, getPosterr := ioutil.ReadAll(getPostResponse.Body)
	if getPosterr != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	var getPostResult map[string]interface{}

	json.Unmarshal([]byte(getPostdata), &getPostResult)

	//getPostResult["caption"] = "Hello@gmail.com"

	time.Sleep(10)

	if getPostResult["_id"] != postIdCreated || getPostResult["caption"] != addPostDetails["caption"] || getPostResult["image_url"] != addPostDetails["image_url"] || getPostResult["user_id"] != addPostDetails["user_id"] {
		t.Error("Wrong Output")
	}
}
