package handler

import (
	"log/slog"
	"net/http"
	"time"

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

func (h *WebhookHandler) HandleAbacatePay(w http.ResponseWriter, r *http.Request) {
	payload, err := webhook.ParseAndValidateWebhook(r, h.cfg.AbacatepayWebhookSecret)
	if err != nil {
		h.logger.Warn("webhook validation failed", "error", err)
		http.Error(w, `{"error":"invalid webhook"}`, http.StatusUnauthorized)
		return
	}

	exists, err := h.webhookRepo.ExistsByEventID(r.Context(), payload.ID)
	if err != nil {
		h.logger.Error("webhook idempotency check failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		w.WriteHeader(http.StatusOK)
		return
	}

	event := &model.WebhookEvent{
		EventID:   payload.ID,
		Event:     payload.Event,
		RawBody:   string(payload.Data),
		Processed: false,
	}
	if err := h.webhookRepo.Insert(r.Context(), event); err != nil {
		h.logger.Error("webhook event insert failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch payload.Event {
	case "checkout.completed":
		if err := h.handleCheckoutCompleted(r, payload); err != nil {
			h.logger.Error("checkout.completed processing failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.logger.Info("checkout.completed processed", "event_id", payload.ID)
	case "checkout.refunded":
		if err := h.handleCheckoutRefunded(r, payload); err != nil {
			h.logger.Error("checkout.refunded processing failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.logger.Info("checkout.refunded processed", "event_id", payload.ID)
	default:
		h.logger.Warn("unknown webhook event", "event", payload.Event)
	}

	if err := h.webhookRepo.MarkAsProcessed(r.Context(), event.ID); err != nil {
		h.logger.Error("failed to mark webhook as processed", "error", err)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) handleCheckoutCompleted(r *http.Request, payload *webhook.AbacatePayWebhookPayload) error {
	data, err := webhook.ParseCheckoutCompleted(payload.Data)
	if err != nil {
		return err
	}

	payment, err := h.paymentRepo.FindByCheckoutID(r.Context(), data.ID)
	if err != nil {
		return err
	}

	if payment.Status == model.PaymentStatusPaid {
		return nil
	}

	now := time.Now()
	if err := h.paymentRepo.UpdateStatus(r.Context(), payment.ID, model.PaymentStatusPaid, &now); err != nil {
		return err
	}

	if err := h.ticketRepo.MarkAsPaid(r.Context(), payment.TicketIDs, payment.BuyerName, payment.BuyerEmail, payment.ID.Hex()); err != nil {
		return err
	}

	return nil
}

func (h *WebhookHandler) handleCheckoutRefunded(r *http.Request, payload *webhook.AbacatePayWebhookPayload) error {
	data, err := webhook.ParseCheckoutCompleted(payload.Data)
	if err != nil {
		return err
	}

	payment, err := h.paymentRepo.FindByCheckoutID(r.Context(), data.ID)
	if err != nil {
		return err
	}

	if err := h.paymentRepo.UpdateStatus(r.Context(), payment.ID, model.PaymentStatusRefunded, nil); err != nil {
		return err
	}

	if err := h.ticketRepo.UpdateManyStatus(r.Context(), payment.TicketIDs, model.TicketStatusAvailable); err != nil {
		return err
	}

	return nil
}
