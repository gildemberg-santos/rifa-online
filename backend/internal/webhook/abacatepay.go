package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AbacatePayWebhookPayload struct {
	ID      string `json:"id"`
	Event   string `json:"event"`
	Data    json.RawMessage `json:"data"`
}

type CheckoutCompletedData struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	ProductID    string `json:"productId"`
	ExternalID   string `json:"externalId"`
	BuyerName    string `json:"buyerName"`
	BuyerEmail   string `json:"buyerEmail"`
	Amount       int    `json:"amount"`
	PaymentMethod string `json:"paymentMethod"`
}

func VerifySignature(publicKey string, body []byte, signatureHeader string) bool {
	sig, err := hex.DecodeString(signatureHeader)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(publicKey))
	mac.Write(body)
	expected := mac.Sum(nil)

	return subtle.ConstantTimeCompare(sig, expected) == 1
}

func ParseAndValidateWebhook(r *http.Request, webhookSecret string) (*AbacatePayWebhookPayload, error) {
	querySecret := r.URL.Query().Get("webhookSecret")
	if querySecret == "" || querySecret != webhookSecret {
		return nil, fmt.Errorf("invalid or missing webhookSecret")
	}

	signature := r.Header.Get("X-Webhook-Signature")
	if signature == "" {
		return nil, fmt.Errorf("missing X-Webhook-Signature header")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	if !VerifySignature(webhookSecret, body, signature) {
		return nil, fmt.Errorf("invalid HMAC signature")
	}

	var payload AbacatePayWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	if payload.Event == "" {
		return nil, fmt.Errorf("missing event type")
	}

	return &payload, nil
}

func ParseCheckoutCompleted(data json.RawMessage) (*CheckoutCompletedData, error) {
	var eventData CheckoutCompletedData
	if err := json.Unmarshal(data, &eventData); err != nil {
		return nil, fmt.Errorf("failed to parse checkout data: %w", err)
	}
	return &eventData, nil
}
