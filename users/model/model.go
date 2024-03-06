package model

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email,omitempty"`
	Type     string             `json:"type,omitempty"`
	Username string             `json:"username,omitempty"`
}

type Login struct {
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	Type     string `json:"type,omitempty"`
	jwt.StandardClaims
}
