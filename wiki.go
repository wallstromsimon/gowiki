package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

//go:embed html_templates/*.html
var templateFS embed.FS

const wikiDataDir = "wiki_data"

var titlePattern = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

type page struct {
	Title string
	Body  string
}

func (p *page) save() error {
	if !titlePattern.MatchString(p.Title) {
		return os.ErrInvalid
	}

	if err := os.MkdirAll(wikiDataDir, 0o755); err != nil {
		return err
	}

	filename := filepath.Join(wikiDataDir, p.Title+".txt")
	return os.WriteFile(filename, []byte(p.Body), 0o600)
}

func loadPage(title string) (*page, error) {
	if !titlePattern.MatchString(title) {
		return nil, os.ErrInvalid
	}

	filename := filepath.Join(wikiDataDir, title+".txt")
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return new(page{Title: title, Body: string(body)}), nil
}

var templates = template.Must(template.ParseFS(templateFS, "html_templates/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func titleFromRequest(w http.ResponseWriter, r *http.Request) (string, bool) {
	title := r.PathValue("title")
	if !titlePattern.MatchString(title) {
		http.NotFound(w, r)
		return "", false
	}
	return title, true
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, ok := titleFromRequest(w, r)
	if !ok {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, ok := titleFromRequest(w, r)
	if !ok {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = new(page{Title: title})
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, ok := titleFromRequest(w, r)
	if !ok {
		return
	}
	body := r.FormValue("body")
	p := new(page{Title: title, Body: body})
	if err := p.save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /view/{title}", viewHandler)
	mux.HandleFunc("GET /edit/{title}", editHandler)
	mux.HandleFunc("POST /save/{title}", saveHandler)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("listening on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
