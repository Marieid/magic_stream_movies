package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	Genre_id   int
	Genre_name string
}

type Ranking struct {
	Ranking_value int
	Ranking_name  string
}

type Movie struct {
	ID           bson.ObjectID `bson:"_id" json:"_id"`
	Imbd_id      string
	Title        string
	Poster_path  string
	YouTube_id   string
	Genre        []Genre
	Admin_review string
	Ranking      Ranking
}
