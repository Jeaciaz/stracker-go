package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/stracker-go/web"
)

func main() {

	// serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// serve dynamic routes
	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/tabs", tabsHandler)
	web.HandleCategoryRoutes()
	web.HandleSpendingRoutes()

	log.Print("Starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/home.html")
	if err != nil {
		log.Fatal("Failed to parse home template", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("Failed to execute home template", err)
	}
}

func tabsHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/tabs.html")
	if err != nil {
		log.Fatal("Failed to parse tabs template", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("Failed to execute tabs template", err)
	}
}
