package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/rifa-online/internal/config"
	"github.com/user/rifa-online/internal/service"
)

func TestAuthHandler_Register_InvalidBody(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, cfg)
	handler := NewAuthHandler(authService)

	req := makeRequest("POST", "/api/v1/auth/register", `{invalid json}`)
	resp := httptest.NewRecorder()

	handler.Register(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, cfg)
	handler := NewAuthHandler(authService)

	req := makeRequest("POST", "/api/v1/auth/login", `{invalid json}`)
	resp := httptest.NewRecorder()

	handler.Login(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestWriteJSON(t *testing.T) {
	resp := httptest.NewRecorder()
	data := map[string]string{"key": "value"}
	writeJSON(resp, http.StatusOK, data)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
	if ct := resp.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestWriteError(t *testing.T) {
	resp := httptest.NewRecorder()
	writeError(resp, "not found", http.StatusNotFound)

	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.Code)
	}
	var body map[string]string
	parseResponse(resp.Result(), &body)
	if body["error"] != "not found" {
		t.Errorf("expected 'not found', got '%s'", body["error"])
	}
}
