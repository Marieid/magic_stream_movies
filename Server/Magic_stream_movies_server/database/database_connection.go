package database

import(
	"fmt"
	"log"
	"os"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"github.com/joho/godotenv"
)

func DBInstance() *mongo.Client{

	err := godotenv.Load(".env")

	if err != nil{
		log.Println("Warning: unable to find the .env file")
	}

	MongoDB := os.Getenv("MONGODB_URI")

	if MongoDB == "" {
		log.Fatal("Fatal: MONGODB_URI not set!")
	}

	fmt.Println("MONGODB_URI: ", MongoDB)

	clientOptions := options.Client().ApplyURI(MongoDB)

	// Client or error object
	client, err := mongo.Connect(clientOptions)

	if err != nil {
		// Raise error to the calling code
		return nil
	}

	return client
}

// Create clien object based on DB instance

var Client *mongo.Client = DBInstance()

// Method to open the connection to the databse

func OpenCollection(collectionName string) *mongo.Collection {

	err := godotenv.Load(".env")

	if err != nil{
		log.Println("Warning: unable to find the .env file")
	}

	databaseName := os.Getenv("DATABASE_NAME")

	fmt.Println("DATABASE_NAME: ", databaseName)

	collection := Client.Database(databaseName).Collection(collectionName)

	if collection == nil {
		log.Println("Error: Unable to retrieve the collection ")
		return nil
	}

	return collection
}