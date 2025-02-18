package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/yuin/goldmark"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
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

// Blog needs Title, Author, Content
type Author struct {
	Name  string
	Email string
}

type PostData struct {
	Title   string
	Author  Author
	Content string
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

		// To render the code snippets also in Markdown using goldmark package from github
		mdRenderer := goldmark.New(
			goldmark.WithExtensions(
				highlighting.NewHighlighting(
					highlighting.WithStyle("dracula"),
				),
			),
		)
		var buf bytes.Buffer
		err = mdRenderer.Convert([]byte(postMarkdown), &buf)
		if err != nil {
			http.Error(w, "Failed to render markdown", http.StatusInternalServerError)
			return
		}

		// Set to HTML
		w.Header().Set("Content-Type", "text/html")

		// Might make the site slow because it gets rendered every time the file changes. Push it before the return statement above
		tpl, err := template.ParseFiles("post.gohtml")
		if (err != nil) {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}

		err = tpl.Execute(w, PostData{
			Title: "My Blog-goo",
			Author: Author{
				Name:  "Shreya P Rao",
				Email: "shreya@example.com",
			},
			Content: buf.String(),
		})

		if (err != nil) {
			http.Error(w, "Error executing the template", http.StatusInternalServerError)
			return
		}
	}
}

// Going to http://localhost:3030/posts/how-to-make-maggi displays the recipe written down in markdown file of the same name!