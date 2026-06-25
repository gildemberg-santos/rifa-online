package handler

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/user/rifa-online/internal/config"
	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
	"github.com/user/rifa-online/internal/service"
	"github.com/user/rifa-online/internal/webhook"
)

type WebhookHandler struct {
	paymentRepo        *repository.PaymentRepo
	ticketRepo         *repository.TicketRepo
	webhookRepo        *repository.WebhookRepo
	userRepo           *repository.UserRepo
	subscriptionSvc    *service.SubscriptionService
	cfg                *config.Config
	logger             *slog.Logger
}

func NewWebhookHandler(
	paymentRepo *repository.PaymentRepo,
	ticketRepo *repository.TicketRepo,
	webhookRepo *repository.WebhookRepo,
	userRepo *repository.UserRepo,
	subscriptionSvc *service.SubscriptionService,
	cfg *config.Config,
	logger *slog.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		paymentRepo:     paymentRepo,
		ticketRepo:      ticketRepo,
		webhookRepo:     webhookRepo,
		userRepo:        userRepo,
		subscriptionSvc: subscriptionSvc,
		cfg:             cfg,
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

	if strings.HasPrefix(payload.OrderNSU, "sub_") {
		h.handleSubscriptionWebhook(w, r, payload)
	} else {
		h.handleRaffleWebhook(w, r, payload)
	}
}

func (h *WebhookHandler) handleSubscriptionWebhook(w http.ResponseWriter, r *http.Request, payload *webhook.InfinitePayWebhookPayload) {
	paymentID := strings.TrimPrefix(payload.OrderNSU, "sub_")

	oid, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		h.logger.Error("invalid subscription payment id", "order_nsu", payload.OrderNSU)
		w.WriteHeader(http.StatusOK)
		return
	}

	payment, err := h.paymentRepo.FindByID(r.Context(), oid)
	if err != nil {
		h.logger.Error("subscription payment not found", "payment_id", paymentID, "error", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	if payment.Status == model.PaymentStatusPaid {
		w.WriteHeader(http.StatusOK)
		return
	}

	paymentMethod := model.PaymentMethodPIX
	if payload.CaptureMethod == "credit_card" {
		paymentMethod = model.PaymentMethodCard
	}

	now := time.Now()
	if err := h.paymentRepo.UpdateStatus(r.Context(), payment.ID, model.PaymentStatusPaid, &now); err != nil {
		h.logger.Error("failed to update subscription payment status", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.paymentRepo.UpdateFields(r.Context(), payment.ID, bson.M{
		"invoiceSlug":    payload.InvoiceSlug,
		"transactionNsu": payload.TransactionNSU,
		"paymentMethod":  paymentMethod,
	}); err != nil {
		h.logger.Error("failed to update subscription payment fields", "error", err)
	}

	if err := h.subscriptionSvc.ActivateSubscription(r.Context(), payment); err != nil {
		h.logger.Error("failed to activate subscription", "user_id", payment.UserID.Hex(), "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info("subscription payment completed via webhook",
		"user_id", payment.UserID.Hex(),
		"invoice_slug", payload.InvoiceSlug,
	)

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) handleRaffleWebhook(w http.ResponseWriter, r *http.Request, payload *webhook.InfinitePayWebhookPayload) {
	payment, err := h.paymentRepo.FindByOrderNSU(r.Context(), payload.OrderNSU)
	if err != nil {
		h.logger.Error("payment not found for order_nsu", "order_nsu", payload.OrderNSU, "error", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	if payment.Status == model.PaymentStatusPaid {
		w.WriteHeader(http.StatusOK)
		return
	}

	paymentMethod := model.PaymentMethodPIX
	if payload.CaptureMethod == "credit_card" {
		paymentMethod = model.PaymentMethodCard
	}

	now := time.Now()
	if err := h.paymentRepo.UpdateStatus(r.Context(), payment.ID, model.PaymentStatusPaid, &now); err != nil {
		h.logger.Error("failed to update payment status", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.paymentRepo.UpdateFields(r.Context(), payment.ID, bson.M{
		"invoiceSlug":    payload.InvoiceSlug,
		"transactionNsu": payload.TransactionNSU,
		"paymentMethod":  paymentMethod,
	}); err != nil {
		h.logger.Error("failed to update payment fields", "error", err)
	}

	if err := h.ticketRepo.MarkAsPaid(r.Context(), payment.TicketIDs, payment.BuyerName, payment.BuyerPhone, payment.ID.Hex()); err != nil {
		h.logger.Error("failed to mark tickets as paid", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info("raffle payment completed via webhook",
		"order_nsu", payload.OrderNSU,
		"invoice_slug", payload.InvoiceSlug,
		"amount", payload.PaidAmount,
		"method", payload.CaptureMethod,
	)

	w.WriteHeader(http.StatusOK)
}


