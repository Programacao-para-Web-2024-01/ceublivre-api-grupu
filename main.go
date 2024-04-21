package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Review struct {
	ProductID string `json:"product_id"`
	UserID    string `json:"user_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
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
	http.Handle("/reviews", &reviews)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
