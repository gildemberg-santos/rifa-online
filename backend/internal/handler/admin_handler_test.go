package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestAdminHandler_UsersNilRepo(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Users panics with nil repo as expected")
		}
	}()

	req := makeRequest("GET", "/api/v1/admin/users", "")
	resp := httptest.NewRecorder()

	handler.Users(resp, req)
	t.Error("Users should have panicked with nil repo")
}

func TestAdminHandler_UserDetails_InvalidHex(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)

	req := makeRequest("GET", "/api/v1/admin/users/invalid", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.UserDetails(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAdminHandler_UserDetailsNilRepo(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("UserDetails panics with nil repo as expected")
		}
	}()

	req := makeRequest("GET", "/api/v1/admin/users/507f1f77bcf86cd799439011", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "507f1f77bcf86cd799439011")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.UserDetails(resp, req)
	t.Error("UserDetails should have panicked with nil repo")
}

func TestAdminHandler_UpdateUserSubscription_InvalidBody(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)

	req := makeRequest("PUT", "/api/v1/admin/users/507f1f77bcf86cd799439011/subscription", `{invalid}`)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "507f1f77bcf86cd799439011")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.UpdateUserSubscription(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAdminHandler_UpdateUserSubscription_InvalidHex(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)

	req := makeRequest("PUT", "/api/v1/admin/users/invalid/subscription", `{"action":"activate"}`)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.UpdateUserSubscription(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAdminHandler_UpdateUserSubscription_InvalidAction(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)

	req := makeRequest("PUT", "/api/v1/admin/users/507f1f77bcf86cd799439011/subscription", `{"action":"invalid_action"}`)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "507f1f77bcf86cd799439011")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.UpdateUserSubscription(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAdminHandler_ConfirmPayment_InvalidHex(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)

	req := makeRequest("POST", "/api/v1/admin/confirm-payment/invalid", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.ConfirmPayment(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestAdminHandler_ConfirmPaymentNilRepo(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("ConfirmPayment panics with nil repo as expected")
		}
	}()

	req := makeRequest("POST", "/api/v1/admin/confirm-payment/507f1f77bcf86cd799439011", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "507f1f77bcf86cd799439011")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.ConfirmPayment(resp, req)
	t.Error("ConfirmPayment should have panicked with nil repo")
}

func TestAdminHandler_RafflesNilRepo(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Raffles panics with nil repo as expected")
		}
	}()

	req := makeRequest("GET", "/api/v1/admin/raffles", "")
	resp := httptest.NewRecorder()

	handler.Raffles(resp, req)
	t.Error("Raffles should have panicked with nil repo")
}

func TestAdminHandler_StatsNilRepo(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Stats panics with nil repo as expected")
		}
	}()

	req := makeRequest("GET", "/api/v1/admin/stats", "")
	resp := httptest.NewRecorder()

	handler.Stats(resp, req)
	t.Error("Stats should have panicked with nil repo")
}
