package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/rifa-online/internal/config"
	"github.com/user/rifa-online/internal/mailer"
	"github.com/user/rifa-online/internal/middleware"
	"github.com/user/rifa-online/internal/service"
)

func TestAuthHandler_Register_InvalidBody(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)
	handler := NewAuthHandler(authService, cfg)

	req := makeRequest("POST", "/api/v1/auth/register", `{invalid json}`)
	resp := httptest.NewRecorder()

	handler.Register(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)
	handler := NewAuthHandler(authService, cfg)

	req := makeRequest("POST", "/api/v1/auth/login", `{invalid}`)
	resp := httptest.NewRecorder()

	handler.Login(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAuthHandler_LoginNilRepo(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)
	handler := NewAuthHandler(authService, cfg)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Login panics with nil repo as expected")
		}
	}()

	req := makeRequest("POST", "/api/v1/auth/login", `{"email":"test@test.com","password":"wrong"}`)
	resp := httptest.NewRecorder()

	handler.Login(resp, req)
	t.Error("Login should have panicked with nil repo")
}

func TestAuthHandler_VerifyEmail_InvalidBody(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)
	handler := NewAuthHandler(authService, cfg)

	req := makeRequest("POST", "/api/v1/auth/verify-email", `{invalid}`)
	resp := httptest.NewRecorder()

	handler.VerifyEmail(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAuthHandler_VerifyEmail_EmptyFields(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)
	handler := NewAuthHandler(authService, cfg)

	req := makeRequest("POST", "/api/v1/auth/verify-email", `{"email":"","code":""}`)
	resp := httptest.NewRecorder()

	handler.VerifyEmail(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAuthHandler_ResendCode_InvalidBody(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)
	handler := NewAuthHandler(authService, cfg)

	req := makeRequest("POST", "/api/v1/auth/resend-code", `{invalid}`)
	resp := httptest.NewRecorder()

	handler.ResendCode(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAuthHandler_ResendCode_EmptyEmail(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)
	handler := NewAuthHandler(authService, cfg)

	req := makeRequest("POST", "/api/v1/auth/resend-code", `{"email":""}`)
	resp := httptest.NewRecorder()

	handler.ResendCode(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAuthHandler_GetProfileNilRepo(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)
	handler := NewAuthHandler(authService, cfg)

	defer func() {
		if r := recover(); r != nil {
			t.Log("GetProfile panics with nil repo as expected")
		}
	}()

	req := makeRequest("GET", "/api/v1/auth/profile", "")
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.GetProfile(resp, req)
	t.Error("GetProfile should have panicked with nil repo")
}

func TestAuthHandler_Refresh_MissingToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)
	handler := NewAuthHandler(authService, cfg)

	req := makeRequest("POST", "/api/v1/auth/refresh", `{}`)
	resp := httptest.NewRecorder()

	handler.Refresh(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)
	handler := NewAuthHandler(authService, cfg)

	req := makeRequest("POST", "/api/v1/auth/logout", "")
	resp := httptest.NewRecorder()

	handler.Logout(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
}

func TestAuthHandler_UpdateProfile_InvalidBody(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := service.NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)
	handler := NewAuthHandler(authService, cfg)

	req := makeRequest("PUT", "/api/v1/auth/profile", `{invalid}`)
	resp := httptest.NewRecorder()

	handler.UpdateProfile(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}
