package main // Defines the package as 'main', indicating an executable program

import ( // Start of the import block for external libraries
    "fmt" // Package for formatted I/O (like printing errors)

    "github.com/gin-gonic/gin" // The Gin web framework, used for building the server and handling HTTP requests

    // Custom package import. This package contains the **handler functions** (Controllers)
    // that implement the **business logic**, which will use the MongoDB connection setup in the 'database' package.
    controller "github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/controllers"
) // End of the import block

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

    // Define a GET route for the path "/movies"
    // This route is handled by the GetMovies function from the imported 'controller' package
    // Retrieves a list of all movies by calling the database functions.
    router.GET("/movies", controller.GetMovies())

    // Define a GET route for the path "/movie/:imdb_id"
    // ":imdb_id" is a **path parameter** that captures a value from the URL (e.g., /movie/tt0133093)
    // This route is handled by the GetMovie function from the 'controller' package
    // Retrieves a single movie's details based on its ID by calling the database functions.
    router.GET("/movie/:imdb_id", controller.GetMovie())

    // Start the server and listen for incoming requests on port 8080
    // router.Run() is a blocking call, meaning the program stays here until the server stops
    if err := router.Run(":8080"); err != nil {
        // If the server fails to start (e.g., port is already in use), an error message is printed
        fmt.Println("Filed to start server", err)
    }

}
