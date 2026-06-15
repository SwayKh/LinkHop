package render

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//go:embed templates/*.html static/*
var templateFS embed.FS

var templates map[string]*template.Template

type Data struct {
	User      *UserInfo
	Error     string
	Success   string
	Links     []LinkInfo
	Link      *LinkInfo
	FormEmail string
}

type UserInfo struct {
	ID         int64
	Email      string
	Username   string
	IsVerified bool
}

type LinkInfo struct {
	ID           int64
	OriginalURL  string
	ShortCode    string
	CustomAlias  string
	ClickCount   int
	LastAccessed string
	CreatedAt    string
	UpdatedAt    string
}

func Load() {
	templates = make(map[string]*template.Template)
	entries, err := templateFS.ReadDir("templates")
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "layout.html" {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".html")
		tmpl := template.Must(template.ParseFS(templateFS, "templates/layout.html", "templates/"+entry.Name()))
		templates[name] = tmpl
	}
	log.Println("Templates loaded")
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	data, err := templateFS.ReadFile("static/favicon.svg")
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(data)
}

func Template(w http.ResponseWriter, name string, data *Data) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Template error (%s): %v", name, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
