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
