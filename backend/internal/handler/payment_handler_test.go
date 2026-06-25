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

func TestPaymentHandler_Checkout_InvalidBody(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	req := makeRequest("POST", "/api/v1/raffles/507f1f77bcf86cd799439011/checkout", `{invalid}`)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "507f1f77bcf86cd799439011")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.Checkout(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestPaymentHandler_Checkout_EmptyNumbers(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	req := makeRequest("POST", "/api/v1/raffles/507f1f77bcf86cd799439011/checkout", `{"numbers":[],"buyerName":"John","buyerPhone":"123"}`)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "507f1f77bcf86cd799439011")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.Checkout(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestPaymentHandler_Checkout_MissingBuyerName(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	req := makeRequest("POST", "/api/v1/raffles/507f1f77bcf86cd799439011/checkout", `{"numbers":[1],"buyerName":"","buyerPhone":"123"}`)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "507f1f77bcf86cd799439011")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.Checkout(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestPaymentHandler_Checkout_MissingBuyerPhone(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	req := makeRequest("POST", "/api/v1/raffles/507f1f77bcf86cd799439011/checkout", `{"numbers":[1],"buyerName":"John","buyerPhone":""}`)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "507f1f77bcf86cd799439011")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.Checkout(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestPaymentHandler_ConfirmPayment_InvalidHex(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	req := makeRequest("POST", "/api/v1/payments/invalid/confirm", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.ConfirmPayment(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestPaymentHandler_MyPayments_MissingPhone(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	req := makeRequest("GET", "/api/v1/payments/mine", "")
	resp := httptest.NewRecorder()

	handler.MyPayments(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestPaymentHandler_GetPayment(t *testing.T) {
	paymentSvc := service.NewPaymentService(nil, nil, nil, nil, nil, nil, nil)
	handler := NewPaymentHandler(paymentSvc, nil, nil, nil)

	req := makeRequest("GET", "/api/v1/payments/someid", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "someid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.GetPayment(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.Code)
	}
}

func TestPaymentHandler_MyTickets_MissingPhone(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	req := makeRequest("GET", "/api/v1/payments/tickets", "")
	resp := httptest.NewRecorder()

	handler.MyTickets(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestPaymentHandler_MyPurchases_InvalidUser(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	req := makeRequest("GET", "/api/v1/payments/purchases", "")
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "invalid-hex")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.MyPurchases(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}
