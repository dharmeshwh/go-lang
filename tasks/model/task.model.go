package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"email,omitempty"`
	Description string             `json:"type,omitempty"`
	Priority    string             `json:"username,omitempty"`
	DueDate     time.Time          `json:"due_date,omitempty"`
}

type Claims struct {
	Username string `json:"username"`
	Type     string `json:"type,omitempty"`
	jwt.StandardClaims
}
