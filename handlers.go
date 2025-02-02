// handlers.go
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// adminLogin handles POST requests for admin login.
// It verifies the credentials from the request against environment variables
// and returns a JWT token if successful.
func adminLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	jwtSecret := os.Getenv("JWT_SECRET")

	if creds.Username != adminUsername || creds.Password != adminPassword {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create a new JWT token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // expires in 1 hour
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	// Return the token.
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// createPost handles POST requests to add a new blog post.
// This endpoint is protected by the requireAuth middleware.
func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post.CreatedAt = time.Now()

	collection := client.Database("my-blog").Collection("posts")
	result, err := collection.InsertOne(context.Background(), post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(post)
}

// getPosts handles GET requests to retrieve all blog posts.
func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var posts []Post

	collection := client.Database("my-blog").Collection("posts")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var post Post
		if err := cursor.Decode(&post); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}

// createComment handles POST requests to add a comment to a blog post.
// This endpoint is public (no authentication required).
func createComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	postIDStr := vars["id"]
	postID, err := primitive.ObjectIDFromHex(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	comment.CreatedAt = time.Now()

	collection := client.Database("my-blog").Collection("posts")
	filter := bson.M{"_id": postID}
	update := bson.M{"$push": bson.M{"comments": comment}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comment)
}
