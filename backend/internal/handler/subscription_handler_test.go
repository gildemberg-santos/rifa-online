package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/rifa-online/internal/middleware"
	"github.com/user/rifa-online/internal/service"
)

func TestSubscriptionHandler_CheckoutNilRepo(t *testing.T) {
	subscriptionSvc := service.NewSubscriptionService(nil, nil, nil, nil)
	handler := NewSubscriptionHandler(subscriptionSvc)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Checkout panics with nil repo as expected")
		}
	}()

	req := makeRequest("POST", "/api/v1/subscriptions/checkout", "")
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.Checkout(resp, req)
	t.Error("Checkout should have panicked with nil repo")
}

func TestSubscriptionHandler_StatusNilRepo(t *testing.T) {
	subscriptionSvc := service.NewSubscriptionService(nil, nil, nil, nil)
	handler := NewSubscriptionHandler(subscriptionSvc)

	defer func() {
		if r := recover(); r != nil {
			t.Log("Status panics with nil repo as expected")
		}
	}()

	req := makeRequest("GET", "/api/v1/subscriptions/status", "")
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.Status(resp, req)
	t.Error("Status should have panicked with nil repo")
}

func TestSubscriptionHandler_UpdateInfinitePayHandle_InvalidBody(t *testing.T) {
	handler := NewSubscriptionHandler(nil)

	req := makeRequest("PUT", "/api/v1/subscriptions/infinitepay", `{invalid}`)
	resp := httptest.NewRecorder()

	handler.UpdateInfinitePayHandle(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.Code)
	}
}

func TestSubscriptionHandler_UpdateInfinitePayHandleNilRepo(t *testing.T) {
	subscriptionSvc := service.NewSubscriptionService(nil, nil, nil, nil)
	handler := NewSubscriptionHandler(subscriptionSvc)

	defer func() {
		if r := recover(); r != nil {
			t.Log("UpdateInfinitePayHandle panics with nil repo as expected")
		}
	}()

	req := makeRequest("PUT", "/api/v1/subscriptions/infinitepay", `{"handle":"test_handle"}`)
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "507f1f77bcf86cd799439011")
	req = req.WithContext(ctx)
	resp := httptest.NewRecorder()

	handler.UpdateInfinitePayHandle(resp, req)
	t.Error("UpdateInfinitePayHandle should have panicked with nil repo")
}
