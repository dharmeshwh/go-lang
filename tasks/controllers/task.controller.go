package taskcontroller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbconfig "tasks/configs"
	authmiddleware "tasks/middleware"
	"tasks/model"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection = dbconfig.DbCollection()

func GetMyAllTasks(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(authmiddleware.UserContextKey).(string)

	if !ok {
		http.Error(w, "Failed to get user details", http.StatusInternalServerError)
		return
	}

	cur, err := collection.Find(context.Background(), bson.M{"username": username})

	if err != nil {
		log.Fatal(err)
	}

	var tasks []primitive.M

	for cur.Next(context.Background()) {
		var task bson.M

		err := cur.Decode(&task)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}

	defer cur.Close(context.Background())

	jsonResponse, err := json.Marshal(tasks)

	if err != nil {
		http.Error(w, "Error encoding users to JSON", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResponse)

}

func CreateOneTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task

	_ = json.NewDecoder(r.Body).Decode(&task)

	_, err := collection.InsertOne(context.Background(), task)

	if err != nil {
		log.Fatal(err)
		return
	}

	jsonResponse, err := json.Marshal(task)

	if err != nil {
		http.Error(w, "Error encoding users to JSON", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResponse)
}

func UpdateOne(w http.ResponseWriter, r *http.Request) {
	var task model.Task

	_ = json.NewDecoder(r.Body).Decode(&task)

	params := mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{
		"title":       task.Title,
		"description": task.Description,
		"priority":    task.Priority,
		"due_date":    task.DueDate,
	}}

	_, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	jsonResponse, err := json.Marshal(task)

	if err != nil {
		http.Error(w, "Error encoding users to JSON", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResponse)
}

func DeleteOneTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}

	_, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(params["id"])
}
