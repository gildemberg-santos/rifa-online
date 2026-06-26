package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/user/rifa-online/internal/middleware"
	"github.com/user/rifa-online/internal/repository"
	"github.com/user/rifa-online/internal/service"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
	paymentRepo    *repository.PaymentRepo
	ticketRepo     *repository.TicketRepo
	userRepo       *repository.UserRepo
}

func NewPaymentHandler(paymentService *service.PaymentService, paymentRepo *repository.PaymentRepo, ticketRepo *repository.TicketRepo, userRepo *repository.UserRepo) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		paymentRepo:    paymentRepo,
		ticketRepo:     ticketRepo,
		userRepo:       userRepo,
	}
}

type checkoutRequest struct {
	Numbers    []int  `json:"numbers"`
	BuyerName  string `json:"buyerName"`
	BuyerEmail string `json:"buyerEmail"`
	BuyerPhone string `json:"buyerPhone"`
}

func (h *PaymentHandler) DevCheckout(w http.ResponseWriter, r *http.Request) {
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
	if req.BuyerName == "" || req.BuyerPhone == "" {
		writeError(w, "buyer name and phone are required", http.StatusBadRequest)
		return
	}

	result, err := h.paymentService.CreateDevCheckout(r.Context(), service.CheckoutInput{
		RaffleID:   raffleID,
		Numbers:    req.Numbers,
		BuyerName:  req.BuyerName,
		BuyerEmail: req.BuyerEmail,
		BuyerPhone: req.BuyerPhone,
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
	if req.BuyerName == "" || req.BuyerPhone == "" {
		writeError(w, "buyer name and phone are required", http.StatusBadRequest)
		return
	}

	result, err := h.paymentService.CreateCheckout(r.Context(), service.CheckoutInput{
		RaffleID:   raffleID,
		Numbers:    req.Numbers,
		BuyerName:  req.BuyerName,
		BuyerEmail: req.BuyerEmail,
		BuyerPhone: req.BuyerPhone,
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
	transactionNSU := r.URL.Query().Get("transaction_nsu")
	slug := r.URL.Query().Get("slug")

	payment, err := h.paymentService.ConfirmRafflePayment(r.Context(), paymentID, transactionNSU, slug)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPaymentID):
			writeError(w, "invalid payment id", http.StatusBadRequest)
		case errors.Is(err, service.ErrPaymentNotFound):
			writeError(w, "payment not found", http.StatusNotFound)
		case errors.Is(err, service.ErrPaymentNotPending):
			writeError(w, "payment is not pending", http.StatusBadRequest)
		case errors.Is(err, service.ErrInvalidPaymentType):
			writeError(w, "invalid payment type", http.StatusBadRequest)
		case errors.Is(err, service.ErrPaymentNotConfirmed), errors.Is(err, service.ErrPaymentAmountMismatch):
			// Provider has not confirmed payment (or underpaid): not yet paid.
			writeError(w, "payment not confirmed", http.StatusPaymentRequired)
		default:
			writeError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": string(payment.Status)})
}

func (h *PaymentHandler) MyPayments(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		writeError(w, "phone query parameter is required", http.StatusBadRequest)
		return
	}

	payments, err := h.paymentRepo.FindByBuyerPhone(r.Context(), phone)
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
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		writeError(w, "phone query parameter is required", http.StatusBadRequest)
		return
	}

	tickets, err := h.ticketRepo.FindPaidByPhone(r.Context(), phone)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, tickets)
}

func (h *PaymentHandler) MyPurchases(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		writeError(w, "invalid user", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.FindByID(r.Context(), oid)
	if err != nil {
		writeError(w, "user not found", http.StatusNotFound)
		return
	}

	items, err := h.paymentService.GetMyPurchases(r.Context(), oid, user.Phone)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, items)
}
