package main // Defines the package as 'main', indicating an executable program

import ( // Start of the import block for external libraries
	"fmt" // Package for formatted I/O (like printing errors)

	"github.com/gin-gonic/gin" // The Gin web framework, used for building the server and handling HTTP requests

	"github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/routes"
)

func main() {
	//This is the main function - the entry point of the application

	// Initialize the Gin router with default middleware (Logger and Recovery)
	router := gin.Default()

	// Define a GET route for the path "/hello"
	// When a request hits this endpoint, the anonymous function (handler) is executed
	router.GET("/hello", func(c *gin.Context) {
		// c.String(200, "Hello, Magic_stream_movies!") sends a simple text response
		// 200 is the HTTP Status Code for "OK"
		c.String(200, "Hello, Magic_stream_movies!")
	})

	routes.SetUpUnprotectedRoutes(router)
	routes.SetUpProtectedRoutes(router)

	// Start the server and listen for incoming requests on port 8080
	// router.Run() is a blocking call, meaning the program stays here until the server stops
	if err := router.Run(":8080"); err != nil {
		// If the server fails to start (e.g., port is already in use), an error message is printed
		fmt.Println("Filed to start server", err)
	}

}
