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
		{"empty buyer name", CheckoutInput{RaffleID: "507f1f77bcf86cd799439011", BuyerName: "", BuyerPhone: "11999999999"}, true},
		{"buyer name too long", CheckoutInput{RaffleID: "507f1f77bcf86cd799439011", BuyerName: string(make([]byte, 151)), BuyerPhone: "11999999999"}, true},
		{"buyer phone too short", CheckoutInput{RaffleID: "507f1f77bcf86cd799439011", BuyerName: "John", BuyerPhone: "123"}, true},
		{"buyer phone too long", CheckoutInput{RaffleID: "507f1f77bcf86cd799439011", BuyerName: "John", BuyerPhone: "1234567890123"}, true},
		{"buyer phone 9 digits invalid", CheckoutInput{RaffleID: "507f1f77bcf86cd799439011", BuyerName: "John", BuyerPhone: "119999999"}, true},
		{"buyer phone 12 digits invalid", CheckoutInput{RaffleID: "507f1f77bcf86cd799439011", BuyerName: "John", BuyerPhone: "119999999999"}, true},
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

func TestPaymentService_CreateDevCheckoutValidation(t *testing.T) {
	svc := NewPaymentService(nil, nil, nil, nil, nil, nil, &config.Config{})

	tests := []struct {
		name    string
		input   CheckoutInput
		wantErr bool
	}{
		{"empty buyer name", CheckoutInput{BuyerName: "", BuyerPhone: "11999999999"}, true},
		{"buyer name too long", CheckoutInput{BuyerName: string(make([]byte, 151)), BuyerPhone: "11999999999"}, true},
		{"buyer phone too short", CheckoutInput{BuyerName: "John", BuyerPhone: "123"}, true},
		{"buyer phone too long", CheckoutInput{BuyerName: "John", BuyerPhone: "1234567890123"}, true},
		{"buyer phone 9 digits invalid", CheckoutInput{BuyerName: "John", BuyerPhone: "119999999"}, true},
		{"buyer phone 12 digits invalid", CheckoutInput{BuyerName: "John", BuyerPhone: "119999999999"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateDevCheckout(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDevCheckout() error = %v, wantErr = %v", err, tt.wantErr)
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

func TestPaymentService_ConfirmRafflePaymentValidation(t *testing.T) {
	svc := NewPaymentService(nil, nil, nil, nil, nil, nil, &config.Config{})

	tests := []struct {
		name      string
		paymentID string
		wantErr   bool
		wantErrIs error
	}{
		{"empty payment id", "", true, ErrInvalidPaymentID},
		{"invalid payment id", "not-a-valid-objectid", true, ErrInvalidPaymentID},
		{"short payment id", "abc", true, ErrInvalidPaymentID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.ConfirmRafflePayment(context.Background(), tt.paymentID, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfirmRafflePayment() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErrIs != nil && err != tt.wantErrIs {
				t.Errorf("ConfirmRafflePayment() wrong error: got %v, want %v", err, tt.wantErrIs)
			}
		})
	}
}

func TestPaymentService_FindPendingRafflePaymentsNilRepo(t *testing.T) {
	svc := NewPaymentService(nil, nil, nil, nil, nil, nil, &config.Config{})

	defer func() {
		if r := recover(); r != nil {
			t.Log("FindPendingRafflePayments panics with nil repo as expected")
		}
	}()

	svc.FindPendingRafflePayments(context.Background())
	t.Error("FindPendingRafflePayments should have panicked with nil repo")
}

func TestPaymentService_GetMyPurchasesNilRepo(t *testing.T) {
	svc := NewPaymentService(nil, nil, nil, nil, nil, nil, &config.Config{})

	defer func() {
		if r := recover(); r != nil {
			t.Log("GetMyPurchases panics with nil repo as expected")
		}
	}()

	svc.GetMyPurchases(context.Background(), [12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, "")
	t.Error("GetMyPurchases should have panicked with nil repo")
}
