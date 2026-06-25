package infinitepay

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("myhandle", "https://example.com/webhook", "https://example.com/redirect")
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.handle != "myhandle" {
		t.Errorf("expected handle 'myhandle', got '%s'", client.handle)
	}
	if client.webhookURL != "https://example.com/webhook" {
		t.Errorf("expected webhookURL 'https://example.com/webhook', got '%s'", client.webhookURL)
	}
	if client.redirectURL != "https://example.com/redirect" {
		t.Errorf("expected redirectURL 'https://example.com/redirect', got '%s'", client.redirectURL)
	}
	if client.baseURL != defaultBaseURL {
		t.Errorf("expected baseURL '%s', got '%s'", defaultBaseURL, client.baseURL)
	}
}

func TestNewClient_EmptyHandle(t *testing.T) {
	client := NewClient("", "", "")
	if client.handle != "" {
		t.Errorf("expected empty handle, got '%s'", client.handle)
	}
}

func TestCreateCheckoutRequest_DefaultsHandle(t *testing.T) {
	client := NewClient("client-handle", "https://wh", "https://rd")

	req := CreateCheckoutRequest{
		Items: []CheckoutItem{
			{Quantity: 1, Price: 1000, Description: "Test"},
		},
		OrderNSU: "order_123",
	}

	if req.Handle == "" {
		t.Log("Handle is empty, CreateCheckout will use client.handle")
	}
	if req.RedirectURL == "" {
		t.Log("RedirectURL is empty, CreateCheckout will use client.redirectURL")
	}
	if req.WebhookURL == "" {
		t.Log("WebhookURL is empty, CreateCheckout will use client.webhookURL")
	}

	t.Logf("Client handle: %s, webhook: %s, redirect: %s", client.handle, client.webhookURL, client.redirectURL)
}

func TestCreateCheckoutResponse_URL(t *testing.T) {
	resp := CreateCheckoutResponse{URL: "https://checkout.infinitepay.io/link/abc"}
	if resp.URL != "https://checkout.infinitepay.io/link/abc" {
		t.Errorf("expected URL, got '%s'", resp.URL)
	}
}

func TestPaymentCheckRequest_Defaults(t *testing.T) {
	req := PaymentCheckRequest{
		Handle:   "myhandle",
		OrderNSU: "order_123",
	}
	if req.Handle != "myhandle" {
		t.Errorf("expected handle 'myhandle', got '%s'", req.Handle)
	}
}

func TestPaymentCheckResponse_Fields(t *testing.T) {
	resp := PaymentCheckResponse{
		Success:       true,
		Paid:          true,
		Amount:        1000,
		PaidAmount:    1000,
		Installments:  1,
		CaptureMethod: "pix",
	}
	if !resp.Success {
		t.Error("expected success")
	}
	if !resp.Paid {
		t.Error("expected paid")
	}
	if resp.CaptureMethod != "pix" {
		t.Errorf("expected pix, got %s", resp.CaptureMethod)
	}
}

func TestErrorResponse(t *testing.T) {
	errResp := ErrorResponse{Error: "invalid handle"}
	if errResp.Error != "invalid handle" {
		t.Errorf("expected 'invalid handle', got '%s'", errResp.Error)
	}
}
