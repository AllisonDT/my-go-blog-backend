// models.go
package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post represents a blog post.
type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	Comments  []Comment          `bson:"comments,omitempty" json:"comments,omitempty"`
}

// Comment represents a comment on a blog post.
type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Author    string             `bson:"author" json:"author"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// Credentials represents login credentials for the admin.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
