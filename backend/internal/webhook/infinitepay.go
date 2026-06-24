package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type InfinitePayWebhookPayload struct {
	InvoiceSlug    string `json:"invoice_slug"`
	Amount         int    `json:"amount"`
	PaidAmount     int    `json:"paid_amount"`
	Installments   int    `json:"installments"`
	CaptureMethod  string `json:"capture_method"`
	TransactionNSU string `json:"transaction_nsu"`
	OrderNSU       string `json:"order_nsu"`
	ReceiptURL     string `json:"receipt_url"`
}

func ParseInfinitePayWebhook(r *http.Request) (*InfinitePayWebhookPayload, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	var payload InfinitePayWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	if payload.OrderNSU == "" {
		return nil, fmt.Errorf("missing order_nsu")
	}

	return &payload, nil
}
