package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/salvaft/link-shortener/cfg"
	"github.com/salvaft/link-shortener/persistance"
	"github.com/salvaft/link-shortener/utils"
	"github.com/salvaft/link-shortener/views"
)

type LinkService struct {
	store persistance.Store

	limiter *Limiter
}

func NewLinkService(store persistance.Store) *LinkService {
	limiter := newLimiter()
	return &LinkService{store, limiter}
}

func (s *LinkService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /", s.limiter.WithLimitsAndValidation(s.handleCreateLinkWeb))
	mux.HandleFunc("POST /api", s.limiter.WithLimitsAndValidation(s.handleCreateLinkAPI))
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

func (s *LinkService) handleCreateLink(w http.ResponseWriter, r *http.Request) (*persistance.Link, string, error) {
	log.Printf("%-20s Request received. Path: %v", "handleCreateLink", r.URL.Path)
	// TODO: Validate url
	href := r.FormValue("href")
	isPresent := true
	url_id, err := s.store.FindURL(href)
	if _, ok := err.(*persistance.URLNotFound); ok {
		// Should create the link
		isPresent = false
	} else if err != nil {
		// Unexpected error
		log.Printf("%-20s Error checking link in db. Err: %v", "handleCreateLink", err)
		w.WriteHeader(http.StatusInternalServerError)
		views.ErrorView("Unexpected Error").Render(r.Context(), w)
		return nil, "", err
	}

	if !isPresent {
		// Creating the link
		log.Printf("%-20s Creating new link. href: %v", "handleCreateLink", href)
		url_id, err = s.store.CreateLink(href)
		if err != nil {
			log.Printf("%-20s Error creating new link. Error: %v", "handleCreateLink", err)
			w.WriteHeader(http.StatusInternalServerError)
			views.ErrorView("Unexpected Error").Render(r.Context(), w)
			return nil, "", err
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
		return nil, "", err
	}
	return &link, signed_token, nil
}

func (s *LinkService) handleCreateLinkAPI(w http.ResponseWriter, r *http.Request) {
	link, signed_token, err := s.handleCreateLink(w, r)
	if err != nil {
		return
	}
	response := struct {
		Signed_token string           `json:"signed_token"`
		Link         persistance.Link `json:"link"`
	}{signed_token, *link}
	b, err := json.Marshal(response)
	if err != nil {
		log.Printf("%-20s Error marshalling link. Error: %v", "handleCreateLinkAPI", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (s *LinkService) handleCreateLinkWeb(w http.ResponseWriter, r *http.Request) {
	link, signed_token, err := s.handleCreateLink(w, r)
	if err != nil {
		return
	}
	err = views.Home(true, signed_token, link).Render(r.Context(), w)
	if err != nil {
		log.Printf("%-20s Error rendering view. Error: %v", "handleCreateLinkWeb", err)
		w.WriteHeader(http.StatusInternalServerError)
		views.ErrorView("Internal Server Error").Render(r.Context(), w)
	}
}
