package handler

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
	"github.com/user/rifa-online/internal/service"
	"github.com/user/rifa-online/internal/webhook"
)

// WebhookHandler receives InfinitePay notifications. The webhook body is NOT
// trusted: it only acts as a trigger to re-verify the payment server-to-server
// via the payment/subscription services (which call InfinitePay CheckPayment).
type WebhookHandler struct {
	paymentService  *service.PaymentService
	subscriptionSvc *service.SubscriptionService
	webhookRepo     *repository.WebhookRepo
	logger          *slog.Logger
}

func NewWebhookHandler(
	paymentService *service.PaymentService,
	subscriptionSvc *service.SubscriptionService,
	webhookRepo *repository.WebhookRepo,
	logger *slog.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		paymentService:  paymentService,
		subscriptionSvc: subscriptionSvc,
		webhookRepo:     webhookRepo,
		logger:          logger,
	}
}

func (h *WebhookHandler) HandleInfinitePay(w http.ResponseWriter, r *http.Request) {
	payload, err := webhook.ParseInfinitePayWebhook(r)
	if err != nil {
		h.logger.Warn("webhook parse failed", "error", err)
		http.Error(w, `{"error":"invalid webhook"}`, http.StatusBadRequest)
		return
	}

	// Idempotency: skip events we already settled. Keyed by transaction NSU
	// (unique per transaction), falling back to order NSU.
	eventID := payload.TransactionNSU
	if eventID == "" {
		eventID = payload.OrderNSU
	}
	if exists, err := h.webhookRepo.ExistsByEventID(r.Context(), eventID); err == nil && exists {
		w.WriteHeader(http.StatusOK)
		return
	}

	if strings.HasPrefix(payload.OrderNSU, "sub_") {
		paymentID := strings.TrimPrefix(payload.OrderNSU, "sub_")
		if err := h.subscriptionSvc.ConfirmSubscriptionPayment(r.Context(), paymentID); err != nil {
			h.logger.Warn("subscription webhook not settled", "order_nsu", payload.OrderNSU, "error", err)
			// Not yet confirmed by provider: don't record the event so a later
			// (genuine) webhook can re-process. Still 200 to ack receipt.
			w.WriteHeader(http.StatusOK)
			return
		}
		h.logger.Info("subscription settled via webhook", "order_nsu", payload.OrderNSU)
	} else {
		slug := payload.InvoiceSlug
		if slug == "" {
			slug = payload.Slug
		}
		if _, err := h.paymentService.ConfirmRafflePayment(r.Context(), payload.OrderNSU, payload.TransactionNSU, slug); err != nil {
			h.logger.Warn("raffle webhook not settled", "order_nsu", payload.OrderNSU, "error", err)
			w.WriteHeader(http.StatusOK)
			return
		}
		h.logger.Info("raffle settled via webhook", "order_nsu", payload.OrderNSU)
	}

	if err := h.webhookRepo.Insert(r.Context(), &model.WebhookEvent{
		EventID:   eventID,
		Event:     "payment.paid",
		Processed: true,
		CreatedAt: time.Now(),
	}); err != nil {
		h.logger.Error("failed to record webhook event", "event_id", eventID, "error", err)
	}

	w.WriteHeader(http.StatusOK)
}
