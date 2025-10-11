package controllers

import (
	"context"  // Package for context handling, crucial for managing request lifecycles and timeouts
	"net/http" // Standard library package for HTTP status codes
	"time"     // Package for managing time and timeouts

	// Custom imports for database connection and data model structure
	"github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/database" // Import the database connection setup
	"github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/models"   // Import the Movie structure definition
	"github.com/Marieid/magic_stream_movies/Server/Magic_stream_movies_server/utils"
	"github.com/gin-gonic/gin" // The Gin web framework
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"  // MongoDB BSON library for query filters
	"go.mongodb.org/mongo-driver/v2/mongo" // MongoDB driver core functionality
	"golang.org/x/crypto/bcrypt"
)

// Hash user password so it will not be inserted as such in the database using golang.org/x/crypto/bcrypt package
func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(HashPassword), nil
}

// userCollection is a global variable holding the handle to the "user" collection in MongoDB.
// It uses the database.OpenCollection function to establish the connection.
var userCollection *mongo.Collection = database.OpenCollection("users")

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		validate := validator.New()

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "detail": err.Error()})
			return
		}

		// Hashed password is stored as the user password
		hashedPassword, err := HashPassword(user.Password)

		// Check for an error during the hashPassword operation (e.g., password could not be hashed)
		if err != nil {
			// Respond with a 500 Internal Server Error if fetching fails
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password", "details": err.Error()})
			return // Stop execution
		}

		// Set a context with a 100-second timeout for the single database query.
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel() // Release context resources on exit

		// Returns the number of documents where the user has the email address contained in the context object or an err if operation fails
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		// Check for an error during the CountDocuments operation (e.g., Can not check if user already has an account in the users collection)
		if err != nil {
			// Respond with a 500 Internal Server Error if fetching fails
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user", "details": err.Error()})
			return // Stop execution
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return // Stop execution
		}
		// Create an user id for the new user
		user.User_ID = bson.NewObjectID().Hex()

		// Include the time when the user profile was created/updated
		user.Created_at = time.Now()
		user.Updated_at = time.Now()

		// Assigns hashed password to the user
		user.Password = hashedPassword

		result, err := userCollection.InsertOne(ctx, user)

		// Check for an error during the insertion operation (e.g., Can not check if user already has an account in the users collection)
		if err != nil {
			// Respond with a 500 Internal Server Error if fetching fails
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create user", "details": err.Error()})
			return // Stop execution
		}

		// Returns a 201 status created and the result
		c.JSON(http.StatusCreated, result)

	}

}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin models.UserLogin

		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		// Set a context with a 100-second timeout for the single database query.
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel() // Release context resources on exit

		var foundUser models.User

		err := userCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&foundUser)

		// Check for an error during the find user operation (e.g., Check if the user is registered with that email address)
		if err != nil {
			// Respond with a 500 Internal Server Error if fetching fails
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password", "details": err.Error()})
			return // Stop execution
		}
		// Compares the hashed password stored in the database with the password given by the user trying to login
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))

		// Check for an error during the compare password operation (e.g., Passwords are not equal)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return // Stop execution
		}

		token, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.First_name, foundUser.Last_name, foundUser.Role, foundUser.User_ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		err = utils.UpdateAllTokens(foundUser.User_ID, token, refreshToken, database.Client)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh tokens"})
			return

		}

		c.JSON(http.StatusOK, models.UserResponse{
			User_ID:          foundUser.User_ID,
			First_name:       foundUser.First_name,
			Last_name:        foundUser.Last_name,
			Email:            foundUser.Email,
			Role:             foundUser.Role,
			Token:            token,
			Refresh_token:    refreshToken,
			Favourite_genres: foundUser.Favourite_genres,
		})

	}
}
