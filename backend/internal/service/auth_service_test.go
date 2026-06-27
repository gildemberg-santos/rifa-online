package service

import (
	"context"
	"strings"
	"testing"

	"github.com/user/rifa-online/internal/config"
	"github.com/user/rifa-online/internal/mailer"
)

func TestAuthService_RegisterValidation(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)

	tests := []struct {
		name    string
		input   RegisterInput
		wantErr bool
	}{
		{"empty name", RegisterInput{Name: "", Email: "a@b.com", Password: "123456"}, true},
		{"empty email", RegisterInput{Name: "Test", Email: "", Password: "123456"}, true},
		{"empty password", RegisterInput{Name: "Test", Email: "a@b.com", Password: ""}, true},
		{"name too short", RegisterInput{Name: "A", Email: "a@b.com", Password: "123456"}, true},
		{"name too long", RegisterInput{Name: strings.Repeat("a", 101), Email: "a@b.com", Password: "123456"}, true},
		{"email too long", RegisterInput{Name: "Test", Email: strings.Repeat("a", 256) + "@b.com", Password: "123456"}, true},
		{"password too short", RegisterInput{Name: "Test", Email: "a@b.com", Password: "12345"}, true},
		{"password too long", RegisterInput{Name: "Test", Email: "a@b.com", Password: strings.Repeat("a", 73)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Register(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_LoginValidation(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)

	tests := []struct {
		name    string
		input   LoginInput
		wantErr bool
	}{
		{"empty email", LoginInput{Email: "", Password: "123"}, true},
		{"empty password", LoginInput{Email: "a@b.com", Password: ""}, true},
		{"empty both", LoginInput{Email: "", Password: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Login(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil && err != ErrInvalidCredentials {
				t.Errorf("Login() wrong error: got %v, want %v", err, ErrInvalidCredentials)
			}
		})
	}
}

func TestAuthService_VerifyEmailValidation(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)

	tests := []struct {
		name    string
		email   string
		code    string
		wantErr bool
	}{
		{"empty email", "", "123456", true},
		{"empty code", "a@b.com", "", true},
		{"empty both", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.VerifyEmail(context.Background(), tt.email, tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyEmail() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_VerifyEmailNilRepo(t *testing.T) {
	svc := NewAuthService(nil, mailer.New("", 0, "", "", ""), nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("VerifyEmail panics with nil repo as expected")
		}
	}()

	svc.VerifyEmail(context.Background(), "nonexistent@test.com", "000000")
	t.Error("VerifyEmail should have panicked with nil repo")
}

func TestAuthService_ResendCodeValidation(t *testing.T) {
	svc := NewAuthService(nil, mailer.New("", 0, "", "", ""), nil)

	err := svc.ResendCode(context.Background(), "")
	if err == nil {
		t.Error("ResendCode should have failed with empty email")
	}
}

func TestAuthService_ResendCodeNilRepo(t *testing.T) {
	svc := NewAuthService(nil, mailer.New("", 0, "", "", ""), nil)

	defer func() {
		if r := recover(); r != nil {
			t.Log("ResendCode panics with nil repo as expected")
		}
	}()

	svc.ResendCode(context.Background(), "nonexistent@test.com")
	t.Error("ResendCode should have panicked with nil repo")
}

func TestAuthService_UpdateProfileValidation(t *testing.T) {
	svc := NewAuthService(nil, mailer.New("", 0, "", "", ""), nil)

	tests := []struct {
		name    string
		userID  string
		input   UpdateProfileInput
		wantErr bool
	}{
		{"invalid user id", "bad-hex", UpdateProfileInput{}, true},
		{"name too short", "507f1f77bcf86cd799439011", UpdateProfileInput{Name: "A"}, true},
		{"name too long", "507f1f77bcf86cd799439011", UpdateProfileInput{Name: strings.Repeat("a", 101)}, true},
		{"email too long", "507f1f77bcf86cd799439011", UpdateProfileInput{Email: strings.Repeat("a", 256) + "@b.com"}, true},
		{"phone too short", "507f1f77bcf86cd799439011", UpdateProfileInput{Phone: "123"}, true},
		{"phone too long", "507f1f77bcf86cd799439011", UpdateProfileInput{Phone: "123456789012"}, true},
		{"password too short", "507f1f77bcf86cd799439011", UpdateProfileInput{Password: "12345"}, true},
		{"password too long", "507f1f77bcf86cd799439011", UpdateProfileInput{Password: strings.Repeat("a", 73)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.UpdateProfile(context.Background(), tt.userID, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateProfile() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_GetProfileInvalidID(t *testing.T) {
	svc := NewAuthService(nil, mailer.New("", 0, "", "", ""), nil)

	_, err := svc.GetProfile(context.Background(), "bad-hex")
	if err != ErrUserNotFound {
		t.Errorf("GetProfile() wrong error: got %v, want %v", err, ErrUserNotFound)
	}
}

func TestAuthService_RefreshTokenInvalid(t *testing.T) {
	svc := NewAuthService(nil, mailer.New("", 0, "", "", ""), &config.Config{JWTSecret: "test-secret"})

	_, err := svc.RefreshToken(context.Background(), "invalid-token")
	if err != ErrInvalidCredentials {
		t.Errorf("RefreshToken() wrong error: got %v, want %v", err, ErrInvalidCredentials)
	}
}

func TestGenerateCodeLength(t *testing.T) {
	code, err := generateCode()
	if err != nil {
		t.Fatalf("generateCode() failed: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("generateCode() length = %d, want 6", len(code))
	}
}

func TestGenerateCodeNumeric(t *testing.T) {
	code, err := generateCode()
	if err != nil {
		t.Fatalf("generateCode() failed: %v", err)
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("generateCode() non-digit character: %c", c)
		}
	}
}

func TestGenerateCodePadding(t *testing.T) {
	for i := 0; i < 100; i++ {
		code, err := generateCode()
		if err != nil {
			t.Fatalf("generateCode() failed: %v", err)
		}
		if len(code) != 6 {
			t.Fatalf("generateCode() length = %d, want 6", len(code))
		}
	}
}

func TestAuthService_NewAuthService(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	mail := mailer.New("", 0, "", "", "")
	svc := NewAuthService(nil, mail, cfg)
	if svc == nil {
		t.Error("NewAuthService returned nil")
	}
	if svc.cfg != cfg {
		t.Error("NewAuthService did not store cfg")
	}
}

func TestAuthService_LoginWhitespaceTrimming(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret"}
	svc := NewAuthService(nil, mailer.New("", 0, "", "", ""), cfg)

	_, err := svc.Login(context.Background(), LoginInput{Email: "  ", Password: "123"})
	if err != ErrInvalidCredentials {
		t.Errorf("Login() should return ErrInvalidCredentials for whitespace-only email, got %v", err)
	}
}
