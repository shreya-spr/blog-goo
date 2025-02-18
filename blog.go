package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
)
func main() {
	mux := http.NewServeMux()

	// HTML template with TailwindCSS styles
	tpl := template.Must(template.ParseFiles("post.gohtml"))

	mux.HandleFunc("GET /posts/{slug}", PostHandler(FileReader{}, tpl))

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
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

type PostData struct {
	Title   string `toml:"title"`
	Author  Author `toml:"author"`
	Content template.HTML
}

// Make HTTP request with POST
func PostHandler(sr SlugReader, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		postMarkdown, err := sr.Read(slug)

		if (err != nil) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		fmt.Println("Pre parse Post markdown:", postMarkdown)
		// FrontMatter
		var post PostData 
		remainingMd, err := frontmatter.Parse(strings.NewReader(postMarkdown), &post)
		fmt.Println("Remaining Markdown:", remainingMd)
		fmt.Println("Post markdown: ", postMarkdown)
		

		if err != nil {
			fmt.Println("Frontmatter Parsing Error:", err)
			http.Error(w, "Error parsing frontmatter", http.StatusInternalServerError)
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
		err = mdRenderer.Convert([]byte(remainingMd), &buf)
		if err != nil {
			http.Error(w, "Failed to render markdown", http.StatusInternalServerError)
			return
		}

		// Set to HTML
		w.Header().Set("Content-Type", "text/html")

		post.Content = template.HTML(buf.String())
		err = tpl.Execute(w, post)

		if (err != nil) {
			http.Error(w, "Error executing the template", http.StatusInternalServerError)
			return
		}
	}
}

// Going to http://localhost:3030/posts/how-to-make-maggi displays the recipe written down in markdown file of the same name!