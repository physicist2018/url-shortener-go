package authenticator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCheckUserCookie(t *testing.T) {
	t.Run("Valid Cookie", func(t *testing.T) {
		userID := uuid.New().String()
		mac := hmac.New(sha256.New, verySecretKey)
		mac.Write([]byte(userID))
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		cookieValue := fmt.Sprintf("%s.%s", userID, signature)
		cookie := &http.Cookie{
			Name:  "user_session",
			Value: cookieValue,
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		w := httptest.NewRecorder()

		returnedUserID, err := checkUserCookie(w, req)
		assert.NoError(t, err)
		assert.Equal(t, userID, returnedUserID)
	})

	t.Run("Invalid Cookie", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "user_session",
			Value: "invalid.cookie.value",
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		w := httptest.NewRecorder()

		_, err := checkUserCookie(w, req)
		assert.Error(t, err)
	})

	t.Run("No Cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		userID, err := checkUserCookie(w, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, userID)
	})
}

func TestAuthMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(domain.UserIDKey{}).(string)
		w.Write([]byte(userID))
	})

	t.Run("Valid Cookie", func(t *testing.T) {
		userID := uuid.New().String()
		mac := hmac.New(sha256.New, verySecretKey)
		mac.Write([]byte(userID))
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		cookieValue := fmt.Sprintf("%s.%s", userID, signature)
		cookie := &http.Cookie{
			Name:  "user_session",
			Value: cookieValue,
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		w := httptest.NewRecorder()

		handler := AuthMiddleware(nextHandler)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, userID, w.Body.String())
	})

	t.Run("Invalid Cookie", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "user_session",
			Value: "invalid.cookie.value",
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		w := httptest.NewRecorder()

		handler := AuthMiddleware(nextHandler)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("No Cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler := AuthMiddleware(nextHandler)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Body.String())
	})
}

func TestAuthMiddlewareFunc(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(domain.UserIDKey{}).(string)
		w.Write([]byte(userID))
	})

	t.Run("Valid Cookie", func(t *testing.T) {
		userID := uuid.New().String()
		mac := hmac.New(sha256.New, verySecretKey)
		mac.Write([]byte(userID))
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		cookieValue := fmt.Sprintf("%s.%s", userID, signature)
		cookie := &http.Cookie{
			Name:  "user_session",
			Value: cookieValue,
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		w := httptest.NewRecorder()

		handler := AuthMiddlewareFunc(nextHandler)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, userID, w.Body.String())
	})

	t.Run("Invalid Cookie", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "user_session",
			Value: "invalid.cookie.value",
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookie)
		w := httptest.NewRecorder()

		handler := AuthMiddlewareFunc(nextHandler)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("No Cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler := AuthMiddlewareFunc(nextHandler)
		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Body.String())
	})
}

func TestCreateUserCookie(t *testing.T) {
	cookie, userID := createUserCookie()

	assert.NotNil(t, cookie)
	assert.NotEmpty(t, userID)
	assert.Equal(t, "user_session", cookie.Name)
	assert.True(t, cookie.Expires.After(time.Now()))
	assert.True(t, cookie.HttpOnly)
	assert.Equal(t, "/", cookie.Path)

	parts := strings.Split(cookie.Value, ".")
	assert.Len(t, parts, 2)
	assert.Equal(t, userID, parts[0])
}

func TestValidateUserCookie(t *testing.T) {
	t.Run("Valid Cookie", func(t *testing.T) {
		userID := uuid.New().String()
		mac := hmac.New(sha256.New, verySecretKey)
		mac.Write([]byte(userID))
		signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		cookieValue := fmt.Sprintf("%s.%s", userID, signature)
		cookie := &http.Cookie{
			Name:  "user_session",
			Value: cookieValue,
		}

		returnedUserID, valid := validateUserCookie(cookie)
		assert.True(t, valid)
		assert.Equal(t, userID, returnedUserID)
	})

	t.Run("Invalid Cookie", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "user_session",
			Value: "invalid.cookie.value",
		}

		_, valid := validateUserCookie(cookie)
		assert.False(t, valid)
	})

	t.Run("No Cookie", func(t *testing.T) {
		_, valid := validateUserCookie(nil)
		assert.False(t, valid)
	})
}
