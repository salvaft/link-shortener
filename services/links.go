package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/salvaft/go-link-shortener/cfg"
	"github.com/salvaft/go-link-shortener/persistance"
	"github.com/salvaft/go-link-shortener/utils"
	"github.com/salvaft/go-link-shortener/views"
	"golang.org/x/time/rate"
)

type LinkService struct {
	store persistance.Store

	limiter *rate.Limiter
}

func NewLinkService(store persistance.Store) *LinkService {
	// 10 per second, burst of 200
	limiter := rate.NewLimiter(rate.Limit(10), 100)
	return &LinkService{store, limiter}
}

func (s *LinkService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /", s.handleCreateLink)
	mux.HandleFunc("GET /{link}/", s.handleGetLink)
}

func (s *LinkService) handleGetLink(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-20s Request received. Path: %v", "handleGetLink", r.URL.Path)
	linkCode := r.PathValue("link")
	decimal_link := utils.Base64ToDecimal(linkCode)
	href, err := s.store.GetLink(decimal_link)

	if err != nil {
		log.Printf("%-20s Error getting link in db. Error: %v", "handleGetLink", err)
		w.WriteHeader(http.StatusNotFound)
		views.ErrorView("Not found").Render(r.Context(), w)
		return
	} else {
		log.Printf("%-20s Redirecting. Link: %v", "handleGetLink", href)
		http.Redirect(w, r, href, http.StatusMovedPermanently)
		return
	}
}

func (s *LinkService) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Allow() {
		w.WriteHeader(http.StatusTooManyRequests)

		views.ErrorView("Too many requests").Render(r.Context(), w)
	}
	log.Printf("%-20s Request received. Path: %v", "handleCreateLink", r.URL.Path)
	originHeader := r.Header.Get("origin")
	// You can also compare it against the Host or X-Forwarded-Host header.
	origin := fmt.Sprintf("http://%s:%s", cfg.GetConfig().Host, cfg.GetConfig().Port)
	if originHeader != origin {
		fmt.Println(originHeader, origin)
		// Invalid request origin
		w.WriteHeader(http.StatusForbidden)

		views.ErrorView("Forbidden").Render(r.Context(), w)
		return
	}
	// Validate CSRF token
	if !utils.ValidateCSRFToken(r) {
		w.WriteHeader(http.StatusForbidden)
		views.ErrorView("Forbidden").Render(r.Context(), w)

		log.Printf("%-20s csrf token not valid", "handleCreateLink")
		return
	}
	// TODO: Validate url
	href := r.FormValue("href")
	var isPresent bool
	url_id, err := s.store.FindURL(href)
	if _, ok := err.(*persistance.URLNotFound); ok {
		// Should create the link
		isPresent = false
	} else if err != nil {
		// Unexpected error
		log.Printf("%-20s Error checking link in db. Err: %v", "handleCreateLink", err)
		w.WriteHeader(http.StatusInternalServerError)
		views.ErrorView("Unexpected Error").Render(r.Context(), w)
		return
	} else {
		isPresent = true
	}
	if !isPresent {
		// Creating the link
		log.Printf("%-20s Creating new link. href: %v", "handleCreateLink", href)
		url_id, err = s.store.CreateLink(href)
		if err != nil {
			log.Printf("%-20s Error creating new link. Error: %v", "handleCreateLink", err)
			w.WriteHeader(http.StatusInternalServerError)
			views.ErrorView("Unexpected Error").Render(r.Context(), w)
			return
		}
	}
	url_code := utils.DecimalToBase64(url_id)
	full_url := fmt.Sprintf("%s/%s", cfg.GetConfig().Host, url_code)
	link := persistance.Link{B64_code: url_code, Href: href, Id: url_id, Url: full_url}
	signed_token, err := utils.SetCSRFToken(w)
	if err != nil {
		log.Printf("%-20s Error generating CSRF token. Error: %v", "handleCreateLink", err)
		w.WriteHeader(http.StatusInternalServerError)
		views.ErrorView("Internal Server Error").Render(r.Context(), w)
		return
	}
	if r.Header.Get("X-From-Js") == "true" {
		// w.Header().Add(, value string)
		response := struct {
			Signed_token string           `json:"signed_token"`
			Link         persistance.Link `json:"link"`
		}{signed_token, link}
		b, err := json.Marshal(response)
		if err != nil {
			log.Printf("%-20s Error marshalling link. Error: %v", "handleCreateLink", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	views.Home(true, signed_token, &link).Render(r.Context(), w)
}
