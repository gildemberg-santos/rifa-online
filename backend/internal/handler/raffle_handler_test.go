package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/user/rifa-online/internal/middleware"
	"github.com/user/rifa-online/internal/service"
)

func TestRaffleHandler_Create_Unauthorized(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("POST", "/api/v1/raffles", `{}`)
	resp := httptest.NewRecorder()

	handler.Create(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestRaffleHandler_Create_InvalidBody(t *testing.T) {
	raffleSvc := service.NewRaffleService(nil, nil, nil, nil)
	handler := NewRaffleHandler(raffleSvc)

	req := makeRequest("POST", "/api/v1/raffles", `{invalid}`)
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.Create(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestRaffleHandler_Create_InvalidDrawDate(t *testing.T) {
	raffleSvc := service.NewRaffleService(nil, nil, nil, nil)
	handler := NewRaffleHandler(raffleSvc)

	req := makeRequest("POST", "/api/v1/raffles", `{"title":"Test","ticketPrice":10,"maxNumbers":100,"drawDate":"invalid-date"}`)
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.Create(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestRaffleHandler_List(t *testing.T) {
	raffleSvc := service.NewRaffleService(nil, nil, nil, nil)
	handler := NewRaffleHandler(raffleSvc)
	_ = handler
}

func TestRaffleHandler_GetDetail_InvalidHex(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("GET", "/api/v1/raffles/invalid", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.GetDetail(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestRaffleHandler_Update_Unauthorized(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("PUT", "/api/v1/raffles/507f1f77bcf86cd799439011", `{}`)
	resp := httptest.NewRecorder()

	handler.Update(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestRaffleHandler_Update_InvalidRaffleID(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("PUT", "/api/v1/raffles/invalid", `{}`)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.Update(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestRaffleHandler_Delete_Unauthorized(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("DELETE", "/api/v1/raffles/507f1f77bcf86cd799439011", "")
	resp := httptest.NewRecorder()

	handler.Delete(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestRaffleHandler_Delete_InvalidRaffleID(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("DELETE", "/api/v1/raffles/invalid", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.Delete(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestRaffleHandler_Cancel_Unauthorized(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("POST", "/api/v1/raffles/507f1f77bcf86cd799439011/cancel", "")
	resp := httptest.NewRecorder()

	handler.Cancel(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestRaffleHandler_Cancel_InvalidRaffleID(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("POST", "/api/v1/raffles/invalid/cancel", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.Cancel(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestRaffleHandler_MyRaffles_Unauthorized(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("GET", "/api/v1/raffles/mine", "")
	resp := httptest.NewRecorder()

	handler.MyRaffles(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestRaffleHandler_Stats_Unauthorized(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("GET", "/api/v1/raffles/507f1f77bcf86cd799439011/stats", "")
	resp := httptest.NewRecorder()

	handler.Stats(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestRaffleHandler_Stats_InvalidRaffleID(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("GET", "/api/v1/raffles/invalid/stats", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.Stats(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestRaffleHandler_DashboardStats_Unauthorized(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("GET", "/api/v1/raffles/dashboard/stats", "")
	resp := httptest.NewRecorder()

	handler.DashboardStats(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestRaffleHandler_Draw_Unauthorized(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("POST", "/api/v1/raffles/507f1f77bcf86cd799439011/draw", "")
	resp := httptest.NewRecorder()

	handler.Draw(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestRaffleHandler_Draw_InvalidRaffleID(t *testing.T) {
	handler := NewRaffleHandler(nil)

	req := makeRequest("POST", "/api/v1/raffles/invalid/draw", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.Draw(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}
