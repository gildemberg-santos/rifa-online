package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/user/rifa-online/internal/middleware"
	"github.com/user/rifa-online/internal/service"
)

type SubscriptionHandler struct {
	subscriptionSvc *service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionSvc *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionSvc: subscriptionSvc}
}

func (h *SubscriptionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	result, err := h.subscriptionSvc.CreateSubscriptionCheckout(r.Context(), userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrSubscriptionAlreadyActive) || errors.Is(err, service.ErrPendingSubscriptionExists) {
			status = http.StatusConflict
		}
		writeError(w, err.Error(), status)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *SubscriptionHandler) Status(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	user, err := h.subscriptionSvc.GetStatus(r.Context(), userID)
	if err != nil {
		writeError(w, "user not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"subscriptionStatus":     user.SubscriptionStatus,
		"subscriptionExpiresAt":  user.SubscriptionExpiresAt,
		"subscriptionIsTrial":    user.SubscriptionIsTrial,
		"hasSubscriptionBefore":  user.HasSubscriptionBefore,
	})
}

type updateInfinitePayHandleRequest struct {
	InfinitePayHandle string `json:"infinitePayHandle"`
}

func (h *SubscriptionHandler) UpdateInfinitePayHandle(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req updateInfinitePayHandleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.subscriptionSvc.UpdateInfinitePayHandle(r.Context(), userID, req.InfinitePayHandle)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			writeError(w, "user not found", http.StatusNotFound)
			return
		}
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"infinitePayHandle": user.InfinitePayHandle,
	})
}
