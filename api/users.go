package api

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)
import . "insta-backend-api/models"

func GetUser(w http.ResponseWriter, r *http.Request, dbHandler BaseHandler) {
	if r.Method == "GET" {
		w.Header().Add("Content-Type", "application/json")

		var splitUserId = strings.Split(r.URL.String(), "/users/")
		if splitUserId[1] != "" {
			userId, _ := primitive.ObjectIDFromHex(splitUserId[1])

			var result bson.M
			err := dbHandler.Users.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&result)
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

func AddUser(w http.ResponseWriter, r *http.Request, dbHandler BaseHandler) {
	if r.Method == "POST" {

		var target map[string]interface{}
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &target)

		var user User

		name, ok := target["name"]

		if ok {
			nameString, _ := name.(string)
			user.Name = nameString
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Name Not Provided"))
			return
		}

		email, ok := target["email"]

		if ok {
			emailString, _ := email.(string)
			user.Email = emailString
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Email Not Provided"))
			return
		}

		password, ok := target["password"]

		if ok {
			passwordString, _ := password.(string)
			hashedPassword, passwordHashingErr := bcrypt.GenerateFromPassword([]byte(passwordString), 8)

			if passwordHashingErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			user.Password = string(hashedPassword)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Password Not Provided"))
			return
		}

		one, err := dbHandler.Users.InsertOne(context.TODO(), user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(one)
	} else {
		w.Write([]byte("Method Not Allowed: " + r.Method))
	}
}
