package main

import (
	"log"
	"net/http"

	"github.com/salvaft/link-shortener/cfg"
	"github.com/salvaft/link-shortener/persistance"
	"github.com/salvaft/link-shortener/services"
	"github.com/salvaft/link-shortener/utils"
	"github.com/salvaft/link-shortener/views"
)

// TODO: Add dockerfile
// TODO: Add readme
func main() {
	cfg := cfg.GetConfig()
	mux := http.NewServeMux()
	db := persistance.NewSQLite(cfg.DbName)
	defer db.Close()

	healthService := services.NewHealthService()
	healthService.RegisterRoutes(mux)
	// TODO: Move to a separate file
	mux.HandleFunc("GET /", handleIndex)
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	store := persistance.NewStorage(db)
	linkService := services.NewLinkService(store)
	linkService.RegisterRoutes(mux)
	utils.Logger.Printf("%-20s Server running on http://%s:%s", "main", cfg.Host, cfg.Port)
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatalf("%-20s Error starting server. Error: %v", "main", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	signed_token, err := utils.SetCSRFToken(w)
	if err != nil {
		utils.Logger.Printf("%-20s Error generating CSRF token. Error: %v", "handleIndex", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	views.Home(false, signed_token, nil).Render(r.Context(), w)
}
