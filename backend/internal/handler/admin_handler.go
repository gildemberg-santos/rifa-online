package handler

import (
	"net/http"

	"github.com/user/rifa-online/internal/repository"
)

type AdminHandler struct {
	userRepo   *repository.UserRepo
	raffleRepo *repository.RaffleRepo
	ticketRepo *repository.TicketRepo
	paymentRepo *repository.PaymentRepo
}

func NewAdminHandler(
	userRepo *repository.UserRepo,
	raffleRepo *repository.RaffleRepo,
	ticketRepo *repository.TicketRepo,
	paymentRepo *repository.PaymentRepo,
) *AdminHandler {
	return &AdminHandler{
		userRepo:    userRepo,
		raffleRepo:  raffleRepo,
		ticketRepo:  ticketRepo,
		paymentRepo: paymentRepo,
	}
}

func (h *AdminHandler) Users(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.FindAll(r.Context())
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *AdminHandler) Raffles(w http.ResponseWriter, r *http.Request) {
	raffles, err := h.raffleRepo.FindAll(r.Context())
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, raffles)
}

type adminStats struct {
	TotalUsers        int   `json:"totalUsers"`
	ActiveUsers       int   `json:"activeUsers"`
	TotalRaffles      int   `json:"totalRaffles"`
	ActiveRaffles     int   `json:"activeRaffles"`
	TotalPaidTickets  int64 `json:"totalPaidTickets"`
	TotalRevenue      int64 `json:"totalRevenue"`
	TrialUsers        int   `json:"trialUsers"`
	PastDueUsers      int   `json:"pastDueUsers"`
}

func (h *AdminHandler) Stats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	totalUsers, _ := h.userRepo.CountAll(ctx)
	activeUsers, _ := h.userRepo.CountBySubscription(ctx, "ACTIVE")
	trialUsers, _ := h.userRepo.CountByTrial(ctx)
	pastDueUsers, _ := h.userRepo.CountBySubscription(ctx, "PAST_DUE")

	totalRaffles, _ := h.raffleRepo.CountAll(ctx)
	activeRaffles, _ := h.raffleRepo.CountByStatus(ctx, "ACTIVE")

	totalPaidTickets, _ := h.ticketRepo.CountAllPaid(ctx)
	totalRevenue, _ := h.paymentRepo.SumAllPaid(ctx)

	writeJSON(w, http.StatusOK, adminStats{
		TotalUsers:       totalUsers,
		ActiveUsers:      activeUsers,
		TotalRaffles:     totalRaffles,
		ActiveRaffles:    activeRaffles,
		TotalPaidTickets: totalPaidTickets,
		TotalRevenue:     totalRevenue,
		TrialUsers:       trialUsers,
		PastDueUsers:     pastDueUsers,
	})
}
