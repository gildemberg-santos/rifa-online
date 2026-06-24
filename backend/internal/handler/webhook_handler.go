package handler

import (
	"log/slog"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/user/rifa-online/internal/config"
	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
	"github.com/user/rifa-online/internal/webhook"
)

type WebhookHandler struct {
	paymentRepo *repository.PaymentRepo
	ticketRepo  *repository.TicketRepo
	webhookRepo *repository.WebhookRepo
	cfg         *config.Config
	logger      *slog.Logger
}

func NewWebhookHandler(
	paymentRepo *repository.PaymentRepo,
	ticketRepo *repository.TicketRepo,
	webhookRepo *repository.WebhookRepo,
	cfg *config.Config,
	logger *slog.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		paymentRepo: paymentRepo,
		ticketRepo:  ticketRepo,
		webhookRepo: webhookRepo,
		cfg:         cfg,
		logger:      logger,
	}
}

func (h *WebhookHandler) HandleInfinitePay(w http.ResponseWriter, r *http.Request) {
	payload, err := webhook.ParseInfinitePayWebhook(r)
	if err != nil {
		h.logger.Warn("webhook parse failed", "error", err)
		http.Error(w, `{"error":"invalid webhook"}`, http.StatusBadRequest)
		return
	}

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

	payment.InvoiceSlug = payload.InvoiceSlug
	payment.TransactionNSU = payload.TransactionNSU
	payment.PaymentMethod = paymentMethod

	now := time.Now()
	if err := h.paymentRepo.UpdateStatus(r.Context(), payment.ID, model.PaymentStatusPaid, &now); err != nil {
		h.logger.Error("failed to update payment status", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.paymentRepo.UpdateFields(r.Context(), payment.ID, bsonUpdateFields(payment)); err != nil {
		h.logger.Error("failed to update payment fields", "error", err)
	}

	if err := h.ticketRepo.MarkAsPaid(r.Context(), payment.TicketIDs, payment.BuyerName, payment.BuyerEmail, payment.ID.Hex()); err != nil {
		h.logger.Error("failed to mark tickets as paid", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info("payment completed via webhook",
		"order_nsu", payload.OrderNSU,
		"invoice_slug", payload.InvoiceSlug,
		"amount", payload.PaidAmount,
		"method", payload.CaptureMethod,
	)

	w.WriteHeader(http.StatusOK)
}

func bsonUpdateFields(payment *model.Payment) primitive.M {
	m := primitive.M{
		"invoiceSlug":    payment.InvoiceSlug,
		"transactionNsu": payment.TransactionNSU,
		"paymentMethod":  payment.PaymentMethod,
	}
	return m
}
