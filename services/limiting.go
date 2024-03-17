package services

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/salvaft/go-link-shortener/cfg"
	"github.com/salvaft/go-link-shortener/utils"
	"github.com/salvaft/go-link-shortener/views"
)

type (
	token   string
	Limiter struct {
		ticker   *time.Ticker
		registry map[token]*limits
		mu       sync.Mutex
	}
	limits struct {
		created time.Time
		count   int
	}
)

func newLimiter() *Limiter {
	limiter := Limiter{registry: make(map[token]*limits), ticker: time.NewTicker(time.Second * 10)}
	go limiter.CleanOld()
	return &limiter
}

func (l *Limiter) CheckLimits(t token) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.registry[t]; !ok {
		l.registry[t] = &limits{created: time.Now(), count: 1}
		return true
	} else {
		if l.registry[t].count > 2 {
			return false
		} else {
			l.registry[t].count = l.registry[t].count + 1
			return true
		}
	}
}

func (l *Limiter) CleanOld() {
	for t := range l.ticker.C {
		l.mu.Lock()
		for k, v := range l.registry {
			if t.Sub(v.created) > time.Second*5 {
				log.Printf("%-20s Cleaning token: %v", "CleanOld", k)
				delete(l.registry, k)
			}
		}
		l.mu.Unlock()

	}
}

func (l *Limiter) WithLimitsAndValidation(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		originHeader := r.Header.Get("origin")
		// You can also compare it against the Host or X-Forwarded-Host header.
		origin := fmt.Sprintf("http://%s:%s", cfg.GetConfig().Host, cfg.GetConfig().Port)
		if originHeader != origin {
			// Invalid request origin
			log.Printf("%-20s Invalid request origin", "WithLimitsAndValidation")
			w.WriteHeader(http.StatusForbidden)
			views.ErrorView("Forbidden").Render(r.Context(), w)
			return
		}
		if !utils.ValidateCSRFToken(r) {
			w.WriteHeader(http.StatusForbidden)
			log.Printf("%-20s csrf token not valid", "WithLimitsAndValidation")
			views.ErrorView("Forbidden").Render(r.Context(), w)
			return
		}
		token := token(r.Header.Get("csrf-token"))
		if !l.CheckLimits(token) {
			w.WriteHeader(http.StatusTooManyRequests)
			log.Printf("%-20s too many request with same token", "WithLimitsAndValidation")
			views.ErrorView("Too many requests").Render(r.Context(), w)
			return
		}
		handler(w, r)
	}
}
