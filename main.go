package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Review struct {
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	Comments  []Comment `json:"comments,omitempty"`
}

type Comment struct {
	UserID    string `json:"user_id"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}

type ProductReviews struct {
	Reviews []Review
	lock    sync.Mutex
}

func (pr *ProductReviews) AddReview(review Review) {
	pr.lock.Lock()
	defer pr.lock.Unlock()
	pr.Reviews = append(pr.Reviews, review)
}

func (pr *ProductReviews) AddComment(reviewID, userID, text string) {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	for i, review := range pr.Reviews {
		if review.ProductID == reviewID {
			review.Comments = append(review.Comments, Comment{
				UserID:    userID,
				Text:      text,
				Timestamp: time.Now().Unix(),
			})
			pr.Reviews[i] = review
			break
		}
	}
}

func (pr *ProductReviews) GetComments(reviewID string) []Comment {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	for _, review := range pr.Reviews {
		if review.ProductID == reviewID {
			return review.Comments
		}
	}
	return nil
}

func (pr *ProductReviews) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var review Review
		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pr.AddReview(review)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Review added")
	case "GET":
		pr.lock.Lock()
		defer pr.lock.Unlock()
		reviewsJSON, err := json.Marshal(pr.Reviews)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(reviewsJSON)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
	}
}

func main() {
	var reviews ProductReviews
	http.HandleFunc("/reviews", reviews.ServeHTTP)

	http.HandleFunc("/reviews/comments/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		reviewID := r.URL.Query().Get("review_id")
		var c Comment
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		reviews.AddComment(reviewID, c.UserID, c.Text)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Comment added")
	})

	http.HandleFunc("/reviews/comments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		reviewID := r.URL.Query().Get("review_id")
		comments := reviews.GetComments(reviewID)
		if comments == nil {
			http.Error(w, "No comments found or review does not exist", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(comments)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
