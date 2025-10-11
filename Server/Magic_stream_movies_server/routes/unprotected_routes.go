package routes

import (
	// Custom package import. This package contains the **handler functions** (Controllers)
	// that implement the **business logic**, which will use the MongoDB connection setup in the 'database' package.
	controller "github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/controllers"
	"github.com/gin-gonic/gin" // The Gin web framework
)

func SetUpUnprotectedRoutes(router *gin.Engine) {

	// Define a GET route for the path "/movies"
	// This route is handled by the GetMovies function from the imported 'controller' package
	// Retrieves a list of all movies by calling the database functions.
	router.GET("/movies", controller.GetMovies())

	// Define a POST route for the path "/register"
	// This route is handled by the RegisterUser function from the 'controller' package
	// Adds a user record to the users collection in the database functions.
	router.POST("/register", controller.RegisterUser())

	// Define a POST route for the path "/login"
	// This route is handled by the LoginUser function from the 'controller' package
	// Logins a registered user using tokens to the application
	router.POST("/login", controller.LoginUser())
}
