package services

import (
	"encoding/json"
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
	mux.HandleFunc("GET /{link}/", s.handleGetLink)
}

func (s *LinkService) handleGetLink(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-20s request received. Path: %v", "handleGetLink", r.URL.Path)
	linkCode := r.PathValue("link")
	decimal_link := utils.Base64ToDecimal(linkCode)
	link, err := s.store.GetLink(decimal_link)

	if err != nil {
		log.Printf("%-20s request received. Error: %v", "handleGetLink", err)
		w.WriteHeader(http.StatusNotFound)
		views.ErrorView("Not found").Render(r.Context(), w)
		return
	} else {
		log.Printf("%-20s request received. Link: %v", "handleGetLink", link)
		http.Redirect(w, r, link, http.StatusMovedPermanently)
		return
	}
}


func (s *LinkService) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	log.Printf("%-20s request received. Path: %v", "handleCreateLink", r.URL.Path)
	// Validate CSRF token
	if !utils.ValidateCSRFToken(r) {
		w.WriteHeader(http.StatusForbidden)
		views.ErrorView("Forbidden").Render(r.Context(), w)
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
	// Progressive enhancement we respond with json
	if _, ok := r.Header["X-From-Js"]; ok {
		// TODO: Fix this marshaling
		b, err := json.Marshal(map[string]string{"href": href, "url": link.Url})
		if err != nil {
			log.Printf("%-20s Error marshaling json. Error: %v", "handleCreateLink", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
		// Normal response we have to set a new csrf token
	} else {
		csrfToken, err := utils.SetCSRFToken(w)
		if err != nil {
			log.Printf("%-20s Error generating CSRF token. Error: %v", "handleCreateLink", err)
			w.WriteHeader(http.StatusInternalServerError)
			views.ErrorView("Internal Server Error").Render(r.Context(), w)
			return
		}
		views.Home(true, csrfToken, &link).Render(r.Context(), w)
	}
}
