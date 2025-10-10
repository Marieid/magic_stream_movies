package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID               bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	User_ID          string        `bson:"user_id" json:"user_id"`
	First_name       string        `bson:"first_name" json:"first_name" validate:"required,min=2,max=100"`
	Last_name        string        `bson:"last_name" json:"last_name" validate:"required,min=2,max=100"`
	Email            string        `bson:"email" json:"email" validate:"required,email"`
	Password         string        `bson:"password" json:"password" validate:"required,min=6"`
	Role             string        `bson:"role" json:"role" validate:"oneof=USER ADMIN"`
	Created_at       time.Time     `bson:"created_at" json:"created_at"`
	Updated_at       time.Time     `bson:"updated_at" json:"updated_at"`
	Token            string        `bson:"token" json:"token"`
	Refresh_token    string        `bson:"refresh_token" json:"refresh_token"`
	Favourite_genres []Genre       `bson:"favourite_genres" json:"favourite_genres" validate:"required,dive"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// DTO Data Transfer Object - Transfer data from frontend to backend or between software
// By using DTO we're only exposing the data that needs to be exposed to the client
type UserResponse struct {
	User_ID          string  `json:"user_id"`
	First_name       string  `json:"first_name"`
	Last_name        string  `json:"last_name"`
	Email            string  `json:"email"`
	Role             string  `json:"role"`
	Favourite_genres []Genre `json:"favourite_genres"`
}
