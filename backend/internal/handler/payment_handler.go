package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
	"github.com/user/rifa-online/internal/service"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
	paymentRepo    *repository.PaymentRepo
	ticketRepo     *repository.TicketRepo
}

func NewPaymentHandler(paymentService *service.PaymentService, paymentRepo *repository.PaymentRepo, ticketRepo *repository.TicketRepo) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		paymentRepo:    paymentRepo,
		ticketRepo:     ticketRepo,
	}
}

type checkoutRequest struct {
	Numbers    []int  `json:"numbers"`
	BuyerName  string `json:"buyerName"`
	BuyerEmail string `json:"buyerEmail"`
}

func (h *PaymentHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	raffleID := chi.URLParam(r, "id")

	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Numbers) == 0 {
		writeError(w, "at least one number is required", http.StatusBadRequest)
		return
	}
	if req.BuyerName == "" || req.BuyerEmail == "" {
		writeError(w, "buyer name and email are required", http.StatusBadRequest)
		return
	}

	result, err := h.paymentService.CreateCheckout(r.Context(), service.CheckoutInput{
		RaffleID:   raffleID,
		Numbers:    req.Numbers,
		BuyerName:  req.BuyerName,
		BuyerEmail: req.BuyerEmail,
	})
	if err != nil {
		if errors.Is(err, service.ErrRaffleNotFound) {
			writeError(w, "raffle not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrRaffleNotActive) {
			writeError(w, "raffle is not active", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrNumbersUnavailable) {
			writeError(w, "one or more numbers are unavailable", http.StatusConflict)
			return
		}
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *PaymentHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")

	oid, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		writeError(w, "invalid payment id", http.StatusBadRequest)
		return
	}

	payment, err := h.paymentRepo.FindByID(r.Context(), oid)
	if err != nil {
		writeError(w, "payment not found", http.StatusNotFound)
		return
	}

	if payment.Status != model.PaymentStatusPending {
		writeError(w, "payment is not pending", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if err := h.paymentRepo.UpdateStatus(r.Context(), payment.ID, model.PaymentStatusPaid, &now); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.ticketRepo.MarkAsPaid(r.Context(), payment.TicketIDs, payment.BuyerName, payment.BuyerEmail, payment.ID.Hex()); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "confirmed"})
}

func (h *PaymentHandler) MyPayments(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		writeError(w, "email query parameter is required", http.StatusBadRequest)
		return
	}

	payments, err := h.paymentService.GetMyPayments(r.Context(), email)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, payments)
}

func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")

	payment, err := h.paymentService.GetPaymentByID(r.Context(), paymentID)
	if err != nil {
		writeError(w, "payment not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, payment)
}

func (h *PaymentHandler) MyTickets(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		writeError(w, "email query parameter is required", http.StatusBadRequest)
		return
	}

	tickets, err := h.paymentService.GetMyTickets(r.Context(), email)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, tickets)
}
