package dbconfigs

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://localhost:27017"
const dbName = "test"
const colName = "users"

var collection *mongo.Collection

func init() {
	clientOption := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to Db")

	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection is ready")
}

func GetCollection() *mongo.Collection {
	return collection
}
