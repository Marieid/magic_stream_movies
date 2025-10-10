package database // Defines the package name as 'database'

import ( // Start of the import block for necessary libraries
    "fmt"   // Package for formatted I/O (like printing database names)
    "log"   // Package for logging messages (warnings and fatal errors)
    "os"    // Package for interacting with the operating system (like reading environment variables)

    // MongoDB driver library to connect to and interact with the database
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"

    // Library to load environment variables from a .env file
    "github.com/joho/godotenv"
)

// DBInstance creates and returns a MongoDB client instance
// The function connects to the MongoDB server using the URI from environment variables.
func DBInstance() *mongo.Client {

    // Load environment variables from a .env file in the current directory
    err := godotenv.Load(".env")

    // Check if loading the .env file resulted in an error
    if err != nil {
        log.Println("Warning: unable to find the .env file")
        // Note: The program continues even if the .env file is missing,
        // relying on system environment variables if they are set.
    }

    // Retrieve the MongoDB connection URI from the environment variables
    MongoDB := os.Getenv("MONGODB_URI")

    // Check if the MONGODB_URI environment variable is empty
    if MongoDB == "" {
        // If it's not set, log a fatal error and exit the application
        log.Fatal("Fatal: MONGODB_URI not set!")
    }

    // Print the URI (usually for debugging/confirmation)
    fmt.Println("MONGODB_URI: ", MongoDB)

    // Create a new client options object and apply the retrieved MongoDB URI
    clientOptions := options.Client().ApplyURI(MongoDB)

    // Attempt to connect to the MongoDB server
    // client is the connected object, or err holds the connection error
    client, err := mongo.Connect(clientOptions)

    if err != nil {
        // If connection fails, print the error and return nil (no client)
        log.Fatal("Error connecting to MongoDB:", err)
        return nil
    }

    // Return the successfully connected MongoDB client object
    return client
}

// ---------------------------------------------------------------------
// Create client object based on DB instance
// This line calls DBInstance() once when the program starts
// and stores the connection in a global package variable 'Client'
var Client *mongo.Client = DBInstance()

// ---------------------------------------------------------------------

// OpenCollection opens and returns a specific MongoDB collection
// It takes the name of the desired collection as a string argument.
func OpenCollection(collectionName string) *mongo.Collection {

    // Reload environment variables (re-loading might be redundant if Client is already set,
    // but ensures DATABASE_NAME is available if the first load failed or was skipped)
    err := godotenv.Load(".env")

    if err != nil {
        log.Println("Warning: unable to find the .env file")
    }

    // Retrieve the target database name from environment variables
    databaseName := os.Getenv("DATABASE_NAME")

    fmt.Println("DATABASE_NAME: ", databaseName)

    // Use the global 'Client' to access the specified database and then the specified collection
    collection := Client.Database(databaseName).Collection(collectionName)

    // Check if the collection object was successfully retrieved
    // (Note: MongoDB collections are created on first use, this checks the retrieval process)
    if collection == nil {
        log.Println("Error: Unable to retrieve the collection ")
        return nil
    }

    // Return the handle to the MongoDB collection
    return collection
}