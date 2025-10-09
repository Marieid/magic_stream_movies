package controllers

import (
	"context"
	"net/http"
	"time"
	"github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/database"
	"github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

func GETMovies() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Query will timeout and clear up resources
		ctx, cancel := context.WithTimeout(c, 100*time.Second)

		defer cancel()

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
		}

		c.JSON(http.StatusOK, movies)

	}
}
