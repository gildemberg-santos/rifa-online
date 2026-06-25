package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

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

func TestAdminHandler_Raffles(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)
	_ = handler
}

func TestAdminHandler_Stats(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil, nil)
	_ = handler
}
