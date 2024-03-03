package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	db := NewSQLite("db")
	defer db.Close()
	store := NewStorage(db)
	store.CreateLink("https://www.google.com")
	linkService := NewLinkService(store)
	mux.HandleFunc("/", handleIndex)
	linkService.RegisterRoutes(mux)

	http.ListenAndServe(":8000", mux)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-20s request received", "handleIndex")
}
