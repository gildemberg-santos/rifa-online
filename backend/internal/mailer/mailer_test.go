package mailer

import (
	"testing"
)

func TestNewMailer(t *testing.T) {
	m := New("", 0, "", "", "")
	if m == nil {
		t.Error("NewMailer returned nil")
	}
	if m.Enabled() {
		t.Error("Mailer should not be enabled with empty config")
	}
}

func TestNewMailerEnabled(t *testing.T) {
	m := New("smtp.example.com", 587, "user", "pass", "from@example.com")
	if m == nil {
		t.Error("NewMailer returned nil")
	}
	if !m.Enabled() {
		t.Error("Mailer should be enabled with valid config")
	}
}

func TestSendDisabled(t *testing.T) {
	m := New("", 0, "", "", "")
	err := m.Send("to@test.com", "Subject", "Body")
	if err != nil {
		t.Errorf("Send() on disabled mailer should not fail: %v", err)
	}
}

func TestSendVerificationCodeDisabled(t *testing.T) {
	m := New("", 0, "", "", "")
	err := m.SendVerificationCode("to@test.com", "123456")
	if err != nil {
		t.Errorf("SendVerificationCode() on disabled mailer should not fail: %v", err)
	}
}

func TestEnabled(t *testing.T) {
	m := New("", 0, "", "", "")
	if m.Enabled() {
		t.Error("Expected disabled")
	}

	m2 := New("smtp.example.com", 587, "user", "pass", "from@example.com")
	if !m2.Enabled() {
		t.Error("Expected enabled")
	}

	m3 := New("smtp.example.com", 587, "", "", "from@example.com")
	if m3.Enabled() {
		t.Error("Expected disabled when user is empty")
	}
}
