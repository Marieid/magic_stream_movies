package controllers // Defines the package name as 'controllers', which holds the application's business logic handlers

import ( // Start of the import block
	"context" // Package for context handling, crucial for managing request lifecycles and timeouts
	"errors"
	"log"
	"net/http" // Standard library package for HTTP status codes
	"os"
	"strings"

	// Custom imports for database connection and data model structure
	"github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/database" // Import the database connection setup
	"github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/models"   // Import the Movie structure definition
	"github.com/go-playground/validator/v10"
	"github.com/tmc/langchaingo/llms/openai"

	// Library to load environment variables from a .env file
	"github.com/joho/godotenv"

	// Third-party imports
	"time" // Package for managing time and timeouts

	"github.com/gin-gonic/gin"             // The Gin web framework
	"go.mongodb.org/mongo-driver/v2/bson"  // MongoDB BSON library for query filters
	"go.mongodb.org/mongo-driver/v2/mongo" // MongoDB driver core functionality
) // End of the import block

// movieCollection is a global variable holding the handle to the "movies" collection in MongoDB.
// It uses the database.OpenCollection function to establish the connection.
var movieCollection *mongo.Collection = database.OpenCollection("movies")
var RankingCollection *mongo.Collection = database.OpenCollection("rankings")

// Validator object for data validation
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

func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID required or movies does not exists"})
			return
		}

		var req struct {
			AdminReview string `json:"admin_review"`
		}

		var res struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		sentiment, rankVal, err := GetReviewRanking(req.AdminReview)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking", "detail": err.Error()})
			return

		}

		filter := bson.M{"imdb_id": movieID}

		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		result, err := movieCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie", "detail": err.Error()})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		res.RankingName = sentiment
		res.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, res)

	}
}

func GetReviewRanking(admin_review string) (string, int, error) {
	rankings, err := GetRankings()

	if err != nil {
		return "", 0, err
	}
	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.Ranking_value != 999 {
			sentimentDelimited = sentimentDelimited + ranking.Ranking_name + ","
		}
	}
	//list of ranking words sentiment
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	err = godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: .env file not found")
	}

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")

	if OpenAiApiKey == "" {
		return "", 0, errors.New("could not read OPENAI_API_KEY")
	}

	llm, err := openai.New(openai.WithToken(OpenAiApiKey))

	if err != nil {
		return "", 0, err
	}

	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")

	//Replace the {rankings} placeholder with the list of sentiment names in the rankings collection
	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimited, 1)

	// Response is the llm call with the base prompt + the admin review argument passed to the function
	response, err := llm.Call(context.Background(), base_prompt+admin_review)

	if err != nil {
		return "", 0, err
	}

	rank_value := 0

	for _, ranking := range rankings {
		if ranking.Ranking_name == response {
			rank_value = ranking.Ranking_value
			break
		}
	}
	return response, rank_value, nil
}

// Returns an array of rankings (from the rankings collection) and an error code
func GetRankings() ([]models.Ranking, error) {

	var rankings []models.Ranking

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Return all the documents in the rankings collection and put it in the cursor variable
	cursor, err := RankingCollection.Find(ctx, bson.M{})

	// if an error occurs during the find operation returns an empty array and the error code
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Add whatever it's in the cursor in the rankings variable
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	// Return the rankings array with the documents from the rankings collections in the database
	return rankings, nil

}
