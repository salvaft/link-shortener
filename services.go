package main

import (
	"fmt"
	"log"
	"net/http"
)

type LinkService struct {
	store Store
}

func NewLinkService(store Store) *LinkService {
	return &LinkService{store}
}

func (s *LinkService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{link}/", s.handleGetLink)
}

func (s *LinkService) handleGetLink(w http.ResponseWriter, r *http.Request) {
	linkCode := r.PathValue("link")
	decimal_link := base64ToDecimal(linkCode)
	link, err := s.store.GetLink(decimal_link)

	if err != nil {
		log.Printf("%-20s request received. Error: %v", "handleGetLink", err)
		// TODO: Respond with beautiful error page
		http.Error(w, string(err.Error()), http.StatusNotFound)
	} else {

		log.Printf("%-20s request received. Link: %v", "handleGetLink", link)
		http.Redirect(w, r, link, http.StatusMovedPermanently)
	}
}

func (s *LinkService) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate URL and length
	link := r.FormValue("link")

	isPresent, url_id, err := s.store.FindURL(link)
	if err != nil {
		log.Printf("%-20s error checking link in db. Err: %v", "handleCreateLink", err)
		// return "", errors.New("error checking existing link")

		http.Error(w, string(err.Error()), http.StatusNotFound)
	}
	if isPresent {
		link := decimalToBase64(url_id)
		fmt.Fprint(w, link)
		return
	}
	_, err = s.store.CreateLink(link)
	if err != nil {

		log.Printf("%-20s request received. Error: %v", "handleCreateLink", err)
		// TODO: Respond with beautiful error page
		http.Error(w, string(err.Error()), http.StatusServiceUnavailable)
	}
}
