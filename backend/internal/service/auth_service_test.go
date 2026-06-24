package service

import (
	"context"
	"testing"

	"github.com/user/rifa-online/internal/config"
)

func TestAuthService_RegisterValidation(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	service := NewAuthService(nil, cfg)

	tests := []struct {
		name    string
		input   RegisterInput
		wantErr bool
	}{
		{"empty name", RegisterInput{Name: "", Email: "a@b.com", Password: "123456"}, true},
		{"empty email", RegisterInput{Name: "Test", Email: "", Password: "123456"}, true},
		{"empty password", RegisterInput{Name: "Test", Email: "a@b.com", Password: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Register(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_LoginValidation(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	service := NewAuthService(nil, cfg)

	tests := []struct {
		name    string
		input   LoginInput
		wantErr bool
	}{
		{"empty email", LoginInput{Email: "", Password: "123"}, true},
		{"empty password", LoginInput{Email: "a@b.com", Password: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Login(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
