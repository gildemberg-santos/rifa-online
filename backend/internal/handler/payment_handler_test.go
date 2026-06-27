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

	req := makeRequest("POST", "/api/v1/raffles/507f1f77bcf86cd799439011/checkout", `{"numbers":[],"buyerName":"John","buyerPhone":"11999999999"}`)
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

	req := makeRequest("POST", "/api/v1/raffles/507f1f77bcf86cd799439011/checkout", `{"numbers":[1],"buyerName":"","buyerPhone":"11999999999"}`)
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

func TestPaymentHandler_Checkout_InvalidRaffleID(t *testing.T) {
	paymentSvc := service.NewPaymentService(nil, nil, nil, nil, nil, nil, nil)
	handler := NewPaymentHandler(paymentSvc, nil, nil, nil)

	req := makeRequest("POST", "/api/v1/raffles/abc/checkout", `{"numbers":[1],"buyerName":"John","buyerPhone":"11999999999"}`)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.Checkout(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.Code)
	}
}

func TestPaymentHandler_DevCheckout_InvalidBody(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	req := makeRequest("POST", "/api/v1/raffles/507f1f77bcf86cd799439011/dev-checkout", `{invalid}`)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "507f1f77bcf86cd799439011")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.DevCheckout(resp, req)

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

func TestPaymentHandler_ConfirmPaymentNilRepo(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("ConfirmPayment panics with nil repo as expected")
		}
	}()

	req := makeRequest("POST", "/api/v1/payments/507f1f77bcf86cd799439011/confirm", "")
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "507f1f77bcf86cd799439011")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	resp := httptest.NewRecorder()

	handler.ConfirmPayment(resp, req)
	t.Error("ConfirmPayment should have panicked with nil repo")
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

func TestPaymentHandler_MyPaymentsNilRepo(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("MyPayments panics with nil repo as expected")
		}
	}()

	req := makeRequest("GET", "/api/v1/payments/mine?phone=11999999999", "")
	resp := httptest.NewRecorder()

	handler.MyPayments(resp, req)
	t.Error("MyPayments should have panicked with nil repo")
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

func TestPaymentHandler_MyTicketsNilRepo(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("MyTickets panics with nil repo as expected")
		}
	}()

	req := makeRequest("GET", "/api/v1/payments/tickets?phone=11999999999", "")
	resp := httptest.NewRecorder()

	handler.MyTickets(resp, req)
	t.Error("MyTickets should have panicked with nil repo")
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

func TestPaymentHandler_MyPurchasesNilRepo(t *testing.T) {
	handler := NewPaymentHandler(nil, nil, nil, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("MyPurchases panics with nil repo as expected")
		}
	}()

	req := makeRequest("GET", "/api/v1/payments/purchases", "")
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.MyPurchases(resp, req)
	t.Error("MyPurchases should have panicked with nil repo")
}
