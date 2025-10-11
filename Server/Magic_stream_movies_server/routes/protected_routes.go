package routes

import (
	// Custom package import. This package contains the **handler functions** (Controllers)
	// that implement the **business logic**, which will use the MongoDB connection setup in the 'database' package.
	controller "github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/controllers"
	"github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/middleware"
	"github.com/gin-gonic/gin" // The Gin web framework
)

func SetUpProtectedRoutes(router *gin.Engine) {
	// Excecution of code will abort if the token is not valid i.e. the user is not a valid registered user or they're not logged in
	router.Use(middleware.AuthMiddleware())

	// Protected endpoint
	// Define a GET route for the path "/movie/:imdb_id"
	// ":imdb_id" is a **path parameter** that captures a value from the URL (e.g., /movie/tt0133093)
	// This route is handled by the GetMovie function from the 'controller' package
	// Retrieves a single movie's details based on its ID by calling the database functions.
	router.GET("/movie/:imdb_id", controller.GetMovie())

	// Protected endpoint
	// Define a POST route for the path "/addmovie"
	// This route is handled by the AddMovie function from the 'controller' package
	// Adds a single movie'to the movie collection in the database functions.
	router.POST("/addmovie", controller.AddMovie())
}
