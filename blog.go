package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /posts/{slug}", PostHandler(FileReader{}))

	err := http.ListenAndServe(":3030", mux) // Run server and listen on port 3030

	if (err != nil) {
		fmt.Printf("Server error!")
		log.Fatal(err)
	}
	// Going to "http://localhost:3030/posts/how-to-whatever" Gives: Post: how-to-whatever ON browser
}

// Read blog posts content from the url slug
type SlugReader interface {
	Read(slug string) (string, error)
}

// File reader implementation to read data from
type FileReader struct {

}

// Implementing the interface
func (fr FileReader) Read(slug string) (string, error) {
	// Open file from local storage
	f, err := os.Open(slug + ".md")
	if (err != nil) {
		return "", err // Remember: return everything in `string, err` format
	}

	// Read and close(delay this part a bit with defer)
	defer f.Close()
	b, err := io.ReadAll(f) // ReadAll returns array of bytes (b)
	if (err != nil) {
		return "", err
	}
	return string(b), nil // Typecasting the bytes to string and returning
}

// Make HTTP request with POST
func PostHandler(sr SlugReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		postMarkdown, err := sr.Read(slug)

		if (err != nil) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "%s", postMarkdown)
	}
}

// Going to http://localhost:3030/posts/how-to-make-maggi displays the recipe written down in markdown file of the same name!