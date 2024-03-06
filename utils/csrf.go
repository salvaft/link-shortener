package utils

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func GenerateCSRFToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func ValidateCSRFToken(r *http.Request) bool {
	cookieCsrfToken, err := r.Cookie("csrf-token")
	requestCsrfToken := r.FormValue("csrf-token")

	// fmt.Println(cookieCsrfToken.Value, "middle", requestCsrfToken)
	if err != nil || cookieCsrfToken.Value != requestCsrfToken {
		return false
	}
	return true
}

func SetCSRFToken(w http.ResponseWriter) (string, error) {
	token, err := GenerateCSRFToken()
	if err != nil {
		return "", err
	}
	cookie := http.Cookie{Name: "csrf-token", Value: token, HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode, Path: "/"}
	http.SetCookie(w, &cookie)
	return token, nil
}
