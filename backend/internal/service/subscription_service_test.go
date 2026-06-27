package service

import (
	"context"
	"testing"

	"github.com/user/rifa-online/internal/config"
)

func TestNewSubscriptionService(t *testing.T) {
	cfg := &config.Config{InfinitePayHandle: "test-handle"}
	svc := NewSubscriptionService(nil, nil, nil, cfg)
	if svc == nil {
		t.Error("NewSubscriptionService returned nil")
	}
	if svc.cfg != cfg {
		t.Error("NewSubscriptionService did not store cfg")
	}
}

func TestSubscriptionService_CreateCheckoutValidation(t *testing.T) {
	svc := NewSubscriptionService(nil, nil, nil, &config.Config{})

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{"empty user id", "", true},
		{"invalid user id", "bad-hex", true},
		{"short user id", "abc123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateSubscriptionCheckout(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSubscriptionCheckout() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscriptionService_CheckSubscriptionValidation(t *testing.T) {
	svc := NewSubscriptionService(nil, nil, nil, &config.Config{})

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{"empty user id", "", true},
		{"invalid user id", "bad-hex", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CheckSubscription(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSubscription() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil {
				t.Logf("CheckSubscription() error = %v", err)
			}
		})
	}
}

func TestSubscriptionService_GetStatusValidation(t *testing.T) {
	svc := NewSubscriptionService(nil, nil, nil, &config.Config{})

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{"empty user id", "", true},
		{"invalid user id", "bad-hex", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GetStatus(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStatus() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil {
				t.Logf("GetStatus() error = %v", err)
			}
		})
	}
}

func TestSubscriptionService_UpdateInfinitePayHandleValidation(t *testing.T) {
	svc := NewSubscriptionService(nil, nil, nil, &config.Config{})

	tests := []struct {
		name    string
		userID  string
		handle  string
		wantErr bool
	}{
		{"empty user id", "", "test-handle", true},
		{"invalid user id", "bad-hex", "test-handle", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.UpdateInfinitePayHandle(context.Background(), tt.userID, tt.handle)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateInfinitePayHandle() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil && err != ErrUserNotFound {
				t.Errorf("UpdateInfinitePayHandle() wrong error: got %v, want %v", err, ErrUserNotFound)
			}
		})
	}
}

func TestSubscriptionService_ActivateSubscriptionNilRepo(t *testing.T) {
	svc := NewSubscriptionService(nil, nil, nil, &config.Config{})

	defer func() {
		if r := recover(); r != nil {
			t.Log("ActivateSubscription panics with nil repos as expected")
		}
	}()

	svc.ActivateSubscription(context.Background(), nil)
	t.Error("ActivateSubscription should have panicked with nil repo")
}


