package auth

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGenerateAccessToken(t *testing.T) {
	userID := primitive.NewObjectID()
	secret := "test-secret"

	token, err := GenerateAccessToken(userID, secret)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	userID := primitive.NewObjectID()
	secret := "test-secret"

	token, err := GenerateRefreshToken(userID, secret)
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	userID := primitive.NewObjectID()
	secret := "test-secret"

	token, err := GenerateAccessToken(userID, secret)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.UserID != userID.Hex() {
		t.Fatalf("expected userID %s, got %s", userID.Hex(), claims.UserID)
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	userID := primitive.NewObjectID()

	token, err := GenerateAccessToken(userID, "correct-secret")
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	_, err = ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	_, err := ValidateToken("invalid-token-string", "secret")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestValidateToken_MalformedToken(t *testing.T) {
	_, err := ValidateToken("this.is.not.a.jwt", "secret")
	if err == nil {
		t.Fatal("expected error for malformed token")
	}
}

func TestGenerateAccessToken_Expiry(t *testing.T) {
	userID := primitive.NewObjectID()
	secret := "test-secret"

	token, err := GenerateAccessToken(userID, secret)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}

	expectedExpiry := time.Now().Add(1 * time.Hour)
	if claims.ExpiresAt.Time.After(expectedExpiry.Add(time.Minute)) {
		t.Fatal("access token expiry too far in the future")
	}
}

func TestGenerateRefreshToken_Expiry(t *testing.T) {
	userID := primitive.NewObjectID()
	secret := "test-secret"

	token, err := GenerateRefreshToken(userID, secret)
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}

	expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
	if claims.ExpiresAt.Time.After(expectedExpiry.Add(time.Hour)) {
		t.Fatal("refresh token expiry too far in the future")
	}
}

func TestValidateToken_Subject(t *testing.T) {
	userID := primitive.NewObjectID()
	secret := "test-secret"

	token, err := GenerateAccessToken(userID, secret)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.Subject != userID.Hex() {
		t.Fatalf("expected subject %s, got %s", userID.Hex(), claims.Subject)
	}
}
