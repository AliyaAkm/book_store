package auth

import (
	"encoding/base64"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func NewAuthMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Basic Auth
			token := strings.TrimPrefix(authHeader, "Basic ")
			decoded, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				http.Error(w, "Invalid authorization header", http.StatusBadRequest)
				return
			}

			credentials := strings.SplitN(string(decoded), ":", 2)
			if len(credentials) != 2 {
				http.Error(w, "Invalid credentials format", http.StatusUnauthorized)
				return
			}

			// Проверяем пользователя в БД
			username, password := credentials[0], credentials[1]
			if !validateUser(db, username, password) {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func validateUser(db *gorm.DB, username, password string) bool {
	var count int64
	db.Table("users").Where("username = ? AND password = ?", username, password).Count(&count)
	return count > 0
}
