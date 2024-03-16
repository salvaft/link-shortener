package services

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/salvaft/go-link-shortener/cfg"
)

// TODO: Unit test with mock storage

func TestHealthService_E2E(t *testing.T) {
	config := cfg.GetConfig()
	apiURL := "http://" + config.Host + ":" + config.Port + "/healthzz/"

	response, err := http.Get(apiURL)
	if err != nil {
		t.Errorf("%-20s error while requesting GET at healtzz Err: %v", "TestLinkService_RegisterRoutes_E2E", err)
		return
	}
	defer response.Body.Close()

	// Print the response status code
	if response.Status != "200 OK" {
		t.Errorf("Expected 200 OK but got %s", response.Status)
	}
}

func TestLinkService_RegisterRoutes_E2E(t *testing.T) {
	config := cfg.GetConfig()

	apiURL := "http://" + config.Host + ":" + config.Port + "/create/"

	// Form data
	formData := url.Values{}
	formData.Set("url", "https://www.example.com")

	// Perform the POST request using http.PostForm
	response, err := http.PostForm(apiURL, formData)
	if err != nil {

		t.Errorf("%-20s error while requesting POST at create Err: %v", "TestLinkService_RegisterRoutes_E2E", err)
		return
	}
	defer response.Body.Close()

	if response.Status != "200 OK" {
		t.Errorf("Expected 200 OK but got %s", response.Status)
	}
}
