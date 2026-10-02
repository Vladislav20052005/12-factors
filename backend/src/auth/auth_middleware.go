package auth

import (
	"ad/jwtToken"
	"context"
	"fmt"
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

func GetCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func ValidationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetCookie(r)
		if err != nil || len(token) == 0 {
			fmt.Println(err)
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		userId, err := jwtToken.ValidateToken(token)
		if err != nil {
			http.Error(w, "Authorization header failed to be valid", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
