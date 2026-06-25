package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/rifa-online/internal/auth"
	"github.com/user/rifa-online/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAuth_MissingHeader(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	handler := Auth(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestAuth_InvalidFormat(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	handler := Auth(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	handler := Auth(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestAuth_ValidToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	userID := primitive.NewObjectID()

	token, err := auth.GenerateAccessToken(userID, cfg.JWTSecret)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	var capturedUserID string
	handler := Auth(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
	if capturedUserID != userID.Hex() {
		t.Errorf("expected userID %s, got %s", userID.Hex(), capturedUserID)
	}
}

func TestAuth_BearerCaseInsensitive(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	userID := primitive.NewObjectID()

	token, err := auth.GenerateAccessToken(userID, cfg.JWTSecret)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	var capturedUserID string
	handler := Auth(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID = UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "bearer "+token)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
	if capturedUserID != userID.Hex() {
		t.Errorf("expected userID %s, got %s", userID.Hex(), capturedUserID)
	}
}

func TestUserIDFromContext_Empty(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	userID := UserIDFromContext(req.Context())
	if userID != "" {
		t.Errorf("expected empty, got %s", userID)
	}
}
