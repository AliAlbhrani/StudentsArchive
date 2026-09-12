// Package models contains the data models for the application
package models

type CreatePostRequest struct {
	UserID int
	Body   struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Images  []string `json:"images"`
	}
}
