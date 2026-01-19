package main

import (
	"elo-app/internal/store"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// 1. Init DB
	storage, err := store.NewStorage("./elo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer storage.DB.Close()

	// 2. Load Templates
	tmpl := template.Must(template.ParseGlob("../../templates/*.html"))

	// 3. Inject Dependencies into Handlers (closure or struct method)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Call storage, use logic, render tmpl
	})

	// 4. Start Server
	log.Println("System online at :8080")
	http.ListenAndServe(":8080", nil)
}
