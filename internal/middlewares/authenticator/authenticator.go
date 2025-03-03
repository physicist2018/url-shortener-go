package authenticator

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/physicist2018/url-shortener-go/internal/domain"
)

var verySecretKey = []byte("я памятник себе воздвиг нерукотоврный")

// Вспомогательная функция для проверки куки и получения userID
func checkUserCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("user_session")
	if err != nil || cookie == nil {
		cookie, userID := createUserCookie()
		http.SetCookie(w, cookie)
		return userID, nil
	}

	userID, valid := validateUserCookie(cookie)
	if !valid {
		return "", fmt.Errorf("unauthorized")
	}
	return userID, nil
}

// Мидлварь для авторизации на уровне роутера
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := checkUserCookie(w, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), domain.UserIDKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Мидлварь для авторизации на уровне конкретной ручки
func AuthMiddlewareFunc(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := checkUserCookie(w, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), domain.UserIDKey{}, userID)
		next(w, r.WithContext(ctx))
	})
}

func createUserCookie() (*http.Cookie, string) {
	userID := uuid.New().String()
	mac := hmac.New(sha256.New, verySecretKey)
	mac.Write([]byte(userID))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	cookieValue := fmt.Sprintf("%s.%s", userID, signature)

	cookie := &http.Cookie{
		Name:     "user_session",
		Value:    cookieValue,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}

	return cookie, userID
}

func validateUserCookie(cookie *http.Cookie) (string, bool) {
	if cookie == nil {
		return "", false
	}

	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		return "", false
	}

	userID := parts[0]
	signature := parts[1]

	mac := hmac.New(sha256.New, verySecretKey)
	mac.Write([]byte(userID))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return "", false
	}

	return userID, true
}
