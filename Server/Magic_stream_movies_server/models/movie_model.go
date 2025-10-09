package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	Genre_id   int `bson:"genre_id" json:"genre_id" validate:"required"`
	Genre_name string `bson:"genre_name" json:"genre_name" validate:"required, min=2, max=100"`
}

type Ranking struct {
	Ranking_value int `bson:"ranking_value" json:"ranking_value" validate:"required"`
	Ranking_name  string `bson:"ranking_name" json:"ranking_name" validate:"required"`
	// validate: "oneof= Excellent Great Good Bad Terrible"
}

type Movie struct {
	ID           bson.ObjectID `bson:"_id" json:"_id" validate:"required"`
	Imbd_id      string `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title        string `bson:"title" json:"title" validate:"required, min=2, max=500"`
	Poster_path  string `bson:"poster_path" json:"poster_path" validate:"required, url"`
	YouTube_id   string `bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genre        []Genre `bson:"genre" json:"genre" validate:"required, dive"`
	Admin_review string `bson:"admin_review" json:"admin_review" validate:"required"`
	Ranking      Ranking `bson:"ranking" json:"ranking" validate:"required"`
}
