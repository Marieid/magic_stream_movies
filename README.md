# 🎬 Magic Stream Movies API
This is the backend server for a movie streaming and recommendation platform. It is built with Go using the Gin web framework and uses MongoDB as the persistent data store.

## ✨ Features
RESTful API: Provides endpoints for fetching and managing movie data.

MongoDB Integration: Uses the official Go driver for database operations.

Environment Variables: Securely manages configuration (like database URIs) using a .env file.

JWT Middleware: Includes authentication middleware to protect sensitive routes.

Data Validation: Uses the go-playground/validator package to ensure data integrity on incoming requests.

## ⚙️ Setup and Installation
Prerequisites
You need the following installed on your system:

Go (version 1.18 or higher)

MongoDB (running locally or a connection string to a cloud instance like MongoDB Atlas)

Steps
Clone the Repository:

Bash

git clone https://github.com/Marieid/magic_stream_movies.git
cd magic_stream_movies/Server/Magic_stream_movies_server
Create Environment File:

Create a file named .env in the root of the server directory (magic_stream_movies/Server/Magic_stream_movies_server/) and add your database credentials:

Code snippet

### .env file content
```MONGODB_URI="mongodb+srv://<username>:<password>@<cluster-name>/..."```
```DATABASE_NAME="magic_stream_db"```
```JWT_SECRET_KEY="your_super_secret_key" # Used for token signing```

#### Install Dependencies:

Install the Go dependencies listed in your go.mod file (which includes Gin, MongoDB driver, godotenv, and validator).

Bash

go mod tidy
Run the Server:

Start the server on the defined port (default is 8080).

Bash

go run main.go
The API should now be running at http://localhost:8080.

🚀 API Endpoints
The following routes are exposed by the server:

Method	Path	Description	Access
GET	/hello	Basic test endpoint. Returns "Hello, Magic_stream_movies!".	Public
GET	/movies	Retrieves a list of all movies in the database.	Public
GET	/movie/:imdb_id	Retrieves details for a single movie based on its imdb_id.	Public
POST	/movie/add	[Protected] Adds a new movie document to the collection. Requires a valid JWT in the Authorization header.	Admin/Auth

Export to Sheets
Example Request for Adding a Movie (POST /movie/add)
You would send a JSON body similar to this (assuming your models.Movie struct includes these fields):

JSON

{
    "imdb_id": "tt0133093",
    "title": "The Matrix",
    "year": 1999,
    "role": "admin"
}
The server will respond with a 201 Created status and the MongoDB insertion result if successful, or a 400 Bad Request if validation fails.

📦 Project Structure Overview
File/Directory	Description
main.go	The entry point. Initializes the Gin router and defines all public API routes.
controllers/	Contains the handler functions (GetMovies, GetMovie, AddMovie, etc.). This is the business logic layer.
database/	Contains the database connection logic (DBInstance, OpenCollection). This handles connecting to MongoDB and retrieving collections.
models/	Contains the Go structs (like Movie) that define the data shape for MongoDB and JSON payloads.
middleware/	Contains middleware functions (like auth_middleware.go) for tasks such as JWT validation and access control.
.env	Configuration file for environment variables (database URI, secrets, etc.).

Export to Sheets
🛡️ Authentication
Protected routes require a valid JSON Web Token (JWT) to be passed in the request header:

Header:

Authorization: Bearer <your_jwt_token>
The token is validated by the AuthMiddleware before the request is allowed to proceed to the controller logic.