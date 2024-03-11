package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/salvaft/go-link-shortener/cfg"
	"github.com/salvaft/go-link-shortener/components"
	"github.com/salvaft/go-link-shortener/persistance"
	"github.com/salvaft/go-link-shortener/utils"
	"github.com/salvaft/go-link-shortener/views"
)

type LinkService struct {
	store persistance.Store
}

func NewLinkService(store persistance.Store) *LinkService {
	return &LinkService{store}
}

func (s *LinkService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /", s.handleCreateLink)
	mux.HandleFunc("POST /link", s.handleSavedLink)
	mux.HandleFunc("GET /{link}/", s.handleGetLink)
}

func (s *LinkService) handleGetLink(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-20s Request received. Path: %v", "handleGetLink", r.URL.Path)
	linkCode := r.PathValue("link")
	decimal_link := utils.Base64ToDecimal(linkCode)
	link, err := s.store.GetLink(decimal_link)

	if err != nil {
		log.Printf("%-20s Error getting link in db. Error: %v", "handleGetLink", err)
		w.WriteHeader(http.StatusNotFound)
		views.ErrorView("Not found").Render(r.Context(), w)
		return
	} else {
		log.Printf("%-20s Redirecting. Link: %v", "handleGetLink", link)
		http.Redirect(w, r, link, http.StatusMovedPermanently)
		return
	}
}

func (s *LinkService) handleSavedLink(w http.ResponseWriter, r *http.Request) {
	// Progressive enhacement. Handles the post requests from htmx
	// The post request might come from the form or from the javascript
	// in order to create link elements based on local storage
	// We might use the same endpoint for both casees if checking the headers
	// but the handler gets more complex
	log.Printf("%-20s Request received. Path: %v", "handleSavedLink", r.URL.Path)
	// Validate CSRF token
	if !utils.ValidateCSRFToken(r) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Forbidden")
		return
	}
	href := r.FormValue("href")
	index, err := s.store.FindURL(href)
	if err != nil {
		log.Printf("%-20s URL not found in server. Error: %v", "handleGetLink", err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Link not found")
		return
	}
	linkCode := utils.DecimalToBase64(index)
	link := persistance.Link{
		B64_code: linkCode,
		Href:     href,
		Id:       index,
		Url:      fmt.Sprintf("%s/%s", cfg.InitConfig().Host, linkCode),
	}
	components.LinkEntry(&link).Render(r.Context(), w)
}

func (s *LinkService) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-20s Request received. Path: %v", "handleCreateLink", r.URL.Path)
	// Validate CSRF token
	if !utils.ValidateCSRFToken(r) {
		w.WriteHeader(http.StatusForbidden)
		views.ErrorView("Forbidden").Render(r.Context(), w)

		log.Printf("%-20s csrf token not valid", "handleCreateLink")
		return
	}
	// TODO: vALIDAte url
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
		url_id, err = s.store.CreateLink(href)
		if err != nil {
			log.Printf("%-20s Error creating new link. Error: %v", "handleCreateLink", err)
			w.WriteHeader(http.StatusInternalServerError)
			views.ErrorView("Unexpected Error").Render(r.Context(), w)
			return
		}
	}
	url_code := utils.DecimalToBase64(url_id)
	full_url := fmt.Sprintf("%s/%s", cfg.InitConfig().Host, url_code)
	link := persistance.Link{B64_code: url_code, Href: href, Id: url_id, Url: full_url}
	csrfToken, err := utils.SetCSRFToken(w)
	if err != nil {
		log.Printf("%-20s Error generating CSRF token. Error: %v", "handleCreateLink", err)
		w.WriteHeader(http.StatusInternalServerError)
		views.ErrorView("Internal Server Error").Render(r.Context(), w)
		return
	}
	if r.Header.Get("Hx-Request") != "" {
		// w.Header().Add(, value string)
		components.LinkEntry(&link).Render(r.Context(), w)
		return
	}

	views.Home(true, csrfToken, &link).Render(r.Context(), w)
}
