// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// client is the MongoDB client used across the project.
var client *mongo.Client

func main() {
	// Load environment variables from .env file.
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found; relying on environment variables.")
	}

	// Retrieve MongoDB URI from environment.
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in environment variables")
	}

	// Connect to MongoDB.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	log.Println("Connected to MongoDB Atlas")

	// Set up router.
	router := mux.NewRouter()

	// Home route.
	router.HandleFunc("/", homeHandler).Methods("GET")

	// Admin login route.
	router.HandleFunc("/admin/login", adminLogin).Methods("POST")

	// Protected route for creating posts.
	router.Handle("/posts", requireAuth(http.HandlerFunc(createPost))).Methods("POST")
	// Public route to get posts.
	router.HandleFunc("/posts", getPosts).Methods("GET")

	// Public route for adding comments.
	router.HandleFunc("/posts/{id}/comments", createComment).Methods("POST")

	// Get port from environment or default to 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set up CORS middleware.
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (adjust for production)
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	log.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler(router)))
}

// homeHandler responds to the base URL.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world! This is my advanced Go blog API."))
}
