package controllers // Defines the package name as 'controllers', which holds the application's business logic handlers

import ( // Start of the import block
	"context"  // Package for context handling, crucial for managing request lifecycles and timeouts
	"net/http" // Standard library package for HTTP status codes

	// Custom imports for database connection and data model structure
	"github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/database" // Import the database connection setup
	"github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/models"   // Import the Movie structure definition
	"github.com/go-playground/validator/v10"

	// Third-party imports
	"time" // Package for managing time and timeouts

	"github.com/gin-gonic/gin"             // The Gin web framework
	"go.mongodb.org/mongo-driver/v2/bson"  // MongoDB BSON library for query filters
	"go.mongodb.org/mongo-driver/v2/mongo" // MongoDB driver core functionality
) // End of the import block

// movieCollection is a global variable holding the handle to the "movies" collection in MongoDB.
// It uses the database.OpenCollection function to establish the connection.
var movieCollection *mongo.Collection = database.OpenCollection("movies")

var validate = validator.New()

// GetMovies is the handler function for the GET /movies route.
// It returns a gin.HandlerFunc, which is the signature Gin uses for route handlers.
func GetMovies() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Create a Context with a 100-second timeout for the database query.
		// This ensures the query doesn't hang indefinitely.
		ctx, cancel := context.WithTimeout(c, 100*time.Second)

		// defer cancel() ensures the context resources are released when the function exits.
		// This prevents memory leaks if the handler finishes before the timeout.
		defer cancel()

		var movies []models.Movie // Declare a slice to hold the decoded movie documents

		// Perform the MongoDB Find operation to retrieve ALL documents from the collection.
		// bson.M{} is an empty filter, meaning "find everything."
		cursor, err := movieCollection.Find(ctx, bson.M{})

		// Check for an error during the Find operation (e.g., connection issue)
		if err != nil {
			// Respond with a 500 Internal Server Error if fetching fails
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
			return // Stop execution
		}
		// defer cursor.Close(ctx) ensures the MongoDB cursor is properly closed after processing the results.
		defer cursor.Close(ctx)

		// Decode all documents retrieved by the cursor into the 'movies' slice.
		if err = cursor.All(ctx, &movies); err != nil {
			// Respond with a 500 Internal Server Error if decoding fails
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
			return // Stop execution
		}

		// Respond with a 200 OK status and the list of movies as JSON
		c.JSON(http.StatusOK, movies)

	}
}

// GetMovie is the handler function for the GET /movie/:imdb_id route.
func GetMovie() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Set a context with a 100-second timeout for the single database query.
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel() // Release context resources on exit

		// Retrieve the 'imdb_id' value from the URL path parameters
		movieID := c.Param("imdb_id")

		// Basic validation: check if the path parameter was present
		if movieID == "" {
			// Respond with a 400 Bad Request if the ID is missing
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var movie models.Movie // Declare a variable to hold the single decoded movie document

		// Perform the MongoDB FindOne operation.
		// The filter is bson.M{"imdb_id": movieID}, which looks for a document where the
		// "imdb_id" field matches the value pulled from the URL.
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			// Check if the error is a "no documents found" error
			// (A generic error is handled as Not Found for simplicity here)
			// Respond with a 404 Not Found status
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found."})
			return
		}

		// Respond with a 200 OK status and the single movie object as JSON
		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)

		defer cancel()

		var movie models.Movie

		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		// Insert validated data in the database
		result, err := movieCollection.InsertOne(ctx, movie)

		// If there's an error send a http internalServerError to the client
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}
		
		c.JSON(http.StatusCreated, result)
	}
}
