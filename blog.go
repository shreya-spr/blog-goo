package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /posts/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug") // Get the url path from here
		fmt.Fprintf(w, "Post: %s", slug)
	})

	err := http.ListenAndServe(":3030", mux) // Run server and listen on port 3030

	if (err != nil) {
		fmt.Printf("Server error!")
		log.Fatal(err)
	}
	// Going to "http://localhost:3030/posts/how-to-whatever" Gives: Post: how-to-whatever ON browser
}