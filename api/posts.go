package api

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

import . "AppointyTask/models"

func AddPost(w http.ResponseWriter, r *http.Request, dbHandler BaseHandler) {
	if r.Method == "POST" {

		var target map[string]interface{}
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &target)

		var post Post

		user_id_interface, ok := target["user_id"]

		if ok {
			userId, _ := primitive.ObjectIDFromHex(user_id_interface.(string))
			post.UserId = userId
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("User Id Not Provided"))
			return
		}

		post.PostedTimestamp = time.Now()

		imageUrl, ok := target["image_url"]

		if ok {
			var imageUrlString = imageUrl.(string)
			parsedUrl, err := url.Parse(imageUrlString)
			if err != nil {
				fmt.Println(err)
			}

			post.ImageUrl = *parsedUrl
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Image URL Not Provided"))
			return
		}

		caption, ok := target["caption"]

		if ok {
			post.Caption = caption.(string)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Caption Not Found"))
			return
		}

		fmt.Println(post)

		one, err := dbHandler.Posts.InsertOne(context.TODO(), post)
		if err != nil {
			return
		}

		json.NewEncoder(w).Encode(one)
	} else {
		w.Write([]byte("Method Not Allowed: " + r.Method))
	}
}

func GetAllPostsOfUser(w http.ResponseWriter, r *http.Request, dbHandler BaseHandler) {
	if r.Method == "GET" {
		w.Header().Add("Content-Type", "application/json")

		var pageNo = r.URL.Query().Get("pageno")

		pageNoIntValue, err := strconv.ParseInt(pageNo, 10, 64)
		if err != nil {
			panic(err)
		}

		var pageSize = r.URL.Query().Get("pagesize")

		pageSizeIntValue, err := strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			panic(err)
		}

		userId, _ := primitive.ObjectIDFromHex(strings.Split(r.URL.Path, "/posts/users/")[1])

		//Sort posts in descending order by timestamp {the newest posts first}
		// Pagination
		opts := options.Find().SetSort(bson.D{{"posted_timestamp", -1}}).SetSkip((pageNoIntValue - 1) * pageSizeIntValue).SetLimit(pageSizeIntValue)

		filter := bson.M{"user_id": userId}

		cursor, err := dbHandler.Posts.Find(context.TODO(), filter, opts)

		if err != nil {
			log.Fatal(err)
		}
		var currentUserPosts []bson.M

		if err = cursor.All(context.TODO(), &currentUserPosts); err != nil {
			log.Fatal(err)
		}

		sendEncodedJsonErr := json.NewEncoder(w).Encode(currentUserPosts)
		if sendEncodedJsonErr != nil {
			return
		}

	} else {
		w.Write([]byte("Method Not Allowed: " + r.Method))
	}
}

func GetPost(w http.ResponseWriter, r *http.Request, dbHandler BaseHandler) {
	if r.Method == "GET" {
		w.Header().Add("Content-Type", "application/json")

		var splitPostId = strings.Split(r.URL.String(), "/posts/")
		if splitPostId[1] != "" {
			postId, _ := primitive.ObjectIDFromHex(splitPostId[1])

			var result bson.M
			err := dbHandler.Posts.FindOne(context.TODO(), bson.M{"_id": postId}).Decode(&result)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return
				}
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode(result)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Endpoint"))
		}

	} else {
		w.Write([]byte("Method Not Allowed: " + r.Method))
	}
}
