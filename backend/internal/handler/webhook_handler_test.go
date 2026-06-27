package handler

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookHandler_HandleInfinitePay_InvalidBody(t *testing.T) {
	handler := NewWebhookHandler(nil, nil, nil, slog.New(slog.DiscardHandler))

	req := makeRequest("POST", "/api/v1/webhooks/infinitepay", `{invalid}`)
	resp := httptest.NewRecorder()

	handler.HandleInfinitePay(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestWebhookHandler_HandleInfinitePay_EmptyBody(t *testing.T) {
	handler := NewWebhookHandler(nil, nil, nil, slog.New(slog.DiscardHandler))

	req := makeRequest("POST", "/api/v1/webhooks/infinitepay", ``)
	resp := httptest.NewRecorder()

	handler.HandleInfinitePay(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestWebhookHandler_HandleInfinitePay_MissingOrderNSU(t *testing.T) {
	handler := NewWebhookHandler(nil, nil, nil, slog.New(slog.DiscardHandler))

	req := makeRequest("POST", "/api/v1/webhooks/infinitepay", `{"type":"payment.paid","data":{"id":"evt_123"}}`)
	resp := httptest.NewRecorder()

	handler.HandleInfinitePay(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400 (missing order_nsu), got %d", resp.Code)
	}
}

func TestWebhookHandler_HandleInfinitePayNilRepoRaffle(t *testing.T) {
	handler := NewWebhookHandler(nil, nil, nil, slog.New(slog.DiscardHandler))

	defer func() {
		if r := recover(); r != nil {
			t.Log("HandleInfinitePay panics with nil webhookRepo as expected")
		}
	}()

	req := makeRequest("POST", "/api/v1/webhooks/infinitepay", `{"order_nsu":"507f1f77bcf86cd799439011","transaction_nsu":"txn_123"}`)
	resp := httptest.NewRecorder()

	handler.HandleInfinitePay(resp, req)
	t.Error("HandleInfinitePay should have panicked with nil repo")
}

func TestWebhookHandler_HandleInfinitePayNilRepoSubscription(t *testing.T) {
	handler := NewWebhookHandler(nil, nil, nil, slog.New(slog.DiscardHandler))

	defer func() {
		if r := recover(); r != nil {
			t.Log("HandleInfinitePay panics with nil webhookRepo as expected")
		}
	}()

	req := makeRequest("POST", "/api/v1/webhooks/infinitepay", `{"order_nsu":"sub_507f1f77bcf86cd799439011","transaction_nsu":"txn_456"}`)
	resp := httptest.NewRecorder()

	handler.HandleInfinitePay(resp, req)
	t.Error("HandleInfinitePay should have panicked with nil repo")
}
