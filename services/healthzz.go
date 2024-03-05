package services

import (
	"fmt"
	"net/http"
)

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthzz/", s.handleHealthCheck)
}

func (s *HealthService) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
