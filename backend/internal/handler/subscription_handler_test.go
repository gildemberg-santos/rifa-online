package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/rifa-online/internal/service"
)

func TestSubscriptionHandler_Checkout(t *testing.T) {
	subscriptionSvc := service.NewSubscriptionService(nil, nil, nil, nil)
	handler := NewSubscriptionHandler(subscriptionSvc)

	req := makeRequest("POST", "/api/v1/subscriptions/checkout", "")
	resp := httptest.NewRecorder()

	handler.Checkout(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.Code)
	}
}

func TestSubscriptionHandler_Status(t *testing.T) {
	subscriptionSvc := service.NewSubscriptionService(nil, nil, nil, nil)
	handler := NewSubscriptionHandler(subscriptionSvc)

	req := makeRequest("GET", "/api/v1/subscriptions/status", "")
	resp := httptest.NewRecorder()

	handler.Status(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.Code)
	}
}

func TestSubscriptionHandler_UpdateInfinitePayHandle_InvalidBody(t *testing.T) {
	handler := NewSubscriptionHandler(nil)

	req := makeRequest("PUT", "/api/v1/subscriptions/infinitepay", `{invalid}`)
	resp := httptest.NewRecorder()

	handler.UpdateInfinitePayHandle(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}
