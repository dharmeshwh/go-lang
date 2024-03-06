package usercontroller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	dbconfigs "users/configs"
	authmiddleware "users/middleware"
	"users/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection = dbconfigs.GetCollection()
var signingKey = []byte("your-secret-key")

func GetMyAllUsers(w http.ResponseWriter, r *http.Request) {
	cur, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}

	var users []primitive.M

	for cur.Next(context.Background()) {
		var user bson.M

		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	defer cur.Close(context.Background())

	response, _ := json.Marshal(users)

	w.Header().Set("Content-Type", "application/json")

	w.Write(response)
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(authmiddleware.UserContextKey).(string)

	if !ok {
		http.Error(w, "Failed to get user details", http.StatusInternalServerError)
		return
	}

	var user model.User

	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)

	if err != nil {
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User

	_ = json.NewDecoder(r.Body).Decode(&user)

	_, err := collection.InsertOne(context.Background(), user)

	if err != nil {
		log.Fatal(err)
	}

	response, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")

	w.Write(response)
}

func DeleteAUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}

	_, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	response, _ := json.Marshal(params["id"])

	w.Header().Set("Content-Type", "application/json")

	w.Write(response)
}

func generateToken(username string, usertype string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &model.Claims{
		Username: username,
		Type:     usertype,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds model.Login

	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		log.Fatal(err)
	}

	var user model.User

	err = collection.FindOne(context.Background(), bson.M{"username": creds.Username}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid username", http.StatusUnauthorized)
			return
		}

		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tokenString, err := generateToken(user.Username, user.Type)

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
}
