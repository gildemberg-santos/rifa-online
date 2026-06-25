package webhook

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseInfinitePayWebhook_Valid(t *testing.T) {
	body := `{
		"invoice_slug": "inv_123",
		"amount": 1000,
		"paid_amount": 1000,
		"installments": 1,
		"capture_method": "pix",
		"transaction_nsu": "txn_123",
		"order_nsu": "order_456",
		"receipt_url": "https://example.com/receipt"
	}`

	req := httptest.NewRequest("POST", "/webhook", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	payload, err := ParseInfinitePayWebhook(req)
	if err != nil {
		t.Fatalf("ParseInfinitePayWebhook: %v", err)
	}
	if payload.InvoiceSlug != "inv_123" {
		t.Errorf("expected inv_123, got %s", payload.InvoiceSlug)
	}
	if payload.OrderNSU != "order_456" {
		t.Errorf("expected order_456, got %s", payload.OrderNSU)
	}
	if payload.Amount != 1000 {
		t.Errorf("expected 1000, got %d", payload.Amount)
	}
	if payload.CaptureMethod != "pix" {
		t.Errorf("expected pix, got %s", payload.CaptureMethod)
	}
}

func TestParseInfinitePayWebhook_EmptyBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/webhook", strings.NewReader(""))

	_, err := ParseInfinitePayWebhook(req)
	if err == nil {
		t.Fatal("expected error for empty body")
	}
}

func TestParseInfinitePayWebhook_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/webhook", strings.NewReader("{invalid json}"))

	_, err := ParseInfinitePayWebhook(req)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseInfinitePayWebhook_MissingOrderNSU(t *testing.T) {
	body := `{
		"invoice_slug": "inv_123",
		"amount": 1000
	}`

	req := httptest.NewRequest("POST", "/webhook", strings.NewReader(body))

	_, err := ParseInfinitePayWebhook(req)
	if err == nil {
		t.Fatal("expected error for missing order_nsu")
	}
}

func TestParseInfinitePayWebhook_ExtraFields(t *testing.T) {
	body := `{
		"order_nsu": "ord_999",
		"invoice_slug": "inv_999",
		"amount": 500,
		"paid_amount": 500,
		"capture_method": "credit_card",
		"unknown_field": "should be ignored"
	}`

	req := httptest.NewRequest("POST", "/webhook", strings.NewReader(body))

	payload, err := ParseInfinitePayWebhook(req)
	if err != nil {
		t.Fatalf("ParseInfinitePayWebhook: %v", err)
	}
	if payload.OrderNSU != "ord_999" {
		t.Errorf("expected ord_999, got %s", payload.OrderNSU)
	}
	if payload.CaptureMethod != "credit_card" {
		t.Errorf("expected credit_card, got %s", payload.CaptureMethod)
	}
}
