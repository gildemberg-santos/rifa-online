package service

import (
	"context"
	"testing"

	"github.com/user/rifa-online/internal/config"
)

func TestNewPaymentService(t *testing.T) {
	cfg := &config.Config{InfinitePayHandle: "test-handle"}
	svc := NewPaymentService(nil, nil, nil, nil, nil, nil, cfg)
	if svc == nil {
		t.Error("NewPaymentService returned nil")
	}
	if svc.cfg != cfg {
		t.Error("NewPaymentService did not store cfg")
	}
}

func TestPaymentService_CheckoutValidation(t *testing.T) {
	svc := NewPaymentService(nil, nil, nil, nil, nil, nil, &config.Config{})

	tests := []struct {
		name    string
		input   CheckoutInput
		wantErr bool
	}{
		{"empty raffle id", CheckoutInput{RaffleID: ""}, true},
		{"invalid raffle id", CheckoutInput{RaffleID: "not-a-valid-objectid"}, true},
		{"short raffle id", CheckoutInput{RaffleID: "abc123"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateCheckout(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCheckout() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestPaymentService_GetPaymentByIDValidation(t *testing.T) {
	svc := NewPaymentService(nil, nil, nil, nil, nil, nil, &config.Config{})

	tests := []struct {
		name      string
		paymentID string
		wantErr   bool
	}{
		{"empty payment id", "", true},
		{"invalid payment id", "not-a-valid-objectid", true},
		{"short payment id", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GetPaymentByID(context.Background(), tt.paymentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPaymentByID() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
