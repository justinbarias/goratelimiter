package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinbarias/goratelimiter/ratelimiter"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var posts []Post

func HandleApiRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func main() {
	router := mux.NewRouter()
	posts = append(posts, Post{ID: "1", Title: "My first post", Body: "This is the content of my first post"})
	posts = append(posts, Post{ID: "2", Title: "My second post", Body: "This is the content of my second post"})
	router.HandleFunc("/posts", HandleApiRequest).Methods("GET")
	router.HandleFunc("/anotherposts", HandleApiRequest).Methods("GET")
	wrappedRouter := ratelimiter.NewRateLimiter(router, 10, 10)
	log.Printf("server is listening at localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", wrappedRouter))

}
