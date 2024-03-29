package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/salvaft/link-shortener/cfg"
)

func GenerateCSRFToken() (string, string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", "", err
	}
	rand.Read(tokenBytes)
	csrfToken := base64.StdEncoding.EncodeToString(tokenBytes)
	secret_bytes := []byte(cfg.GetConfig().Secret)
	mac := hmac.New(sha256.New, secret_bytes)
	mac.Write([]byte(csrfToken))
	csrfTokenHMAC := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return csrfToken, csrfTokenHMAC, nil
}

func ValidateCSRFToken(r *http.Request) bool {
	cookieCsrfToken, err := r.Cookie("csrf-token")
	requestCsrfToken := r.FormValue("csrf-token")

	secret_bytes := []byte(cfg.GetConfig().Secret)
	mac := hmac.New(sha256.New, secret_bytes)
	mac.Write([]byte(cookieCsrfToken.Value))
	csrfTokenHMAC := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if err != nil || csrfTokenHMAC != requestCsrfToken {
		return false
	}
	return true
}

func SetCSRFToken(w http.ResponseWriter) (string, error) {
	token, hmac_token, err := GenerateCSRFToken()
	if err != nil {
		return "", err
	}
	cookie := http.Cookie{Name: "csrf-token", Value: token, HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode, Path: "/"}
	http.SetCookie(w, &cookie)
	return hmac_token, nil
}
