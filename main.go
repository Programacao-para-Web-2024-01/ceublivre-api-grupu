package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
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

type Question struct {
    ProductID string    `json:"product_id"`
    UserID    string    `json:"user_id"`
    SellerID  string    `json:"seller_id"`
    Query     string    `json:"query"`
    Timestamp int64     `json:"timestamp"`
    Answers   []Answer  `json:"answers,omitempty"`
}

type Answer struct {
    SellerID  string `json:"seller_id"`
    Response  string `json:"response"`
    Timestamp int64  `json:"timestamp"`
}

type ProductReviews struct {
    Reviews []Review
    lock    sync.Mutex
}

type ProductQuestions struct {
    Questions []Question
    lock      sync.Mutex
}

var bannedWords []string
var lock sync.Mutex

func loadBannedWords(filePath string) error {
    lock.Lock()
    defer lock.Unlock()

    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        bannedWords = append(bannedWords, strings.ToLower(scanner.Text()))
    }

    if err := scanner.Err(); err != nil {
        return err
    }
    return nil
}

func containsBannedWords(text string) bool {
    lock.Lock()
    defer lock.Unlock()

    lowerText := strings.ToLower(text)
    for _, word := range bannedWords {
        if strings.Contains(lowerText, word) {
            return true
        }
    }
    return false
}

func (pr *ProductReviews) AddReview(review Review) error {
    if containsBannedWords(review.Comment) {
        return fmt.Errorf("the review contains banned words")
    }
    pr.lock.Lock()
    defer pr.lock.Unlock()
    pr.Reviews = append(pr.Reviews, review)
    return nil
}

func (pr *ProductReviews) AddComment(reviewID, userID, text string) error {
    if containsBannedWords(text) {
        return fmt.Errorf("the comment contains banned words")
    }
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
    return nil
}

func (pq *ProductQuestions) AddQuestion(question Question) error {
    if containsBannedWords(question.Query) {
        return fmt.Errorf("the question contains banned words")
    }
    pq.lock.Lock()
    defer pq.lock.Unlock()
    pq.Questions = append(pq.Questions, question)
    return nil
}

func (pq *ProductQuestions) AddAnswer(productID, userID, query, sellerID, response string) error {
    if containsBannedWords(response) {
        return fmt.Errorf("the answer contains banned words")
    }
    pq.lock.Lock()
    defer pq.lock.Unlock()
    for i, question := range pq.Questions {
        if question.ProductID == productID && question.UserID == userID && question.Query == query && question.SellerID == sellerID {
            question.Answers = append(question.Answers, Answer{
                SellerID:  sellerID,
                Response:  response,
                Timestamp: time.Now().Unix(),
            })
            pq.Questions[i] = question
            break
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
        if err := pr.AddReview(review); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
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
    if err := loadBannedWords("palavras.txt"); err != nil {
        log.Fatalf("Failed to load banned words: %v", err)
    }

    var reviews ProductReviews
    var questions ProductQuestions

    http.HandleFunc("/reviews", reviews.ServeHTTP)
    http.HandleFunc("/questions/add", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        var question Question
        if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        if err := questions.AddQuestion(question); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        w.WriteHeader(http.StatusCreated)
        fmt.Fprintln(w, "Question added")
    })

    http.HandleFunc("/questions/answer", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        var data struct {
            ProductID string `json:"product_id"`
            UserID    string `json:"user_id"`
            Query     string `json:"query"`
            SellerID  string `json:"seller_id"`
            Response  string `json:"response"`
        }
        if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        if err := questions.AddAnswer(data.ProductID, data.UserID, data.Query, data.SellerID, data.Response); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        w.WriteHeader(http.StatusCreated)
        fmt.Fprintln(w, "Answer added")
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
