package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
)

type AdminHandler struct {
	userRepo    *repository.UserRepo
	raffleRepo  *repository.RaffleRepo
	ticketRepo  *repository.TicketRepo
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

func (h *AdminHandler) UserDetails(w http.ResponseWriter, r *http.Request) {
	userID := 	chi.URLParam(r, "id")
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := h.userRepo.FindByID(ctx, oid)
	if err != nil {
		writeError(w, "user not found", http.StatusNotFound)
		return
	}

	raffles, _ := h.raffleRepo.FindByOrganizer(ctx, oid)
	payments, _ := h.paymentRepo.FindByUserID(ctx, oid)
	tickets, _ := h.ticketRepo.FindPaidByPhone(ctx, user.Phone)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user":     user,
		"raffles":  raffles,
		"payments": payments,
		"tickets":  tickets,
	})
}

type updateSubscriptionRequest struct {
	Action string `json:"action"`
}

func (h *AdminHandler) UpdateUserSubscription(w http.ResponseWriter, r *http.Request) {
	userID := 	chi.URLParam(r, "id")
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req updateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	now := time.Now()

	switch req.Action {
	case "activate":
		expiresAt := now.AddDate(0, 1, 0)
		if err := h.userRepo.UpdateFields(ctx, oid, primitive.M{
			"subscriptionStatus":     model.SubscriptionStatusActive,
			"subscriptionExpiresAt":  expiresAt,
			"subscriptionIsTrial":    false,
			"hasSubscriptionBefore":  true,
		}); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "cancel":
		if err := h.userRepo.UpdateFields(ctx, oid, primitive.M{
			"subscriptionStatus": model.SubscriptionStatusCancelled,
		}); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "past_due":
		if err := h.userRepo.UpdateFields(ctx, oid, primitive.M{
			"subscriptionStatus": model.SubscriptionStatusPastDue,
		}); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "inactive":
		if err := h.userRepo.UpdateFields(ctx, oid, primitive.M{
			"subscriptionStatus":     model.SubscriptionStatusInactive,
			"subscriptionIsTrial":    false,
			"hasSubscriptionBefore":  false,
		}); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		writeError(w, "invalid action: must be activate, cancel, past_due, or inactive", http.StatusBadRequest)
		return
	}

	user, _ := h.userRepo.FindByID(ctx, oid)
	writeJSON(w, http.StatusOK, user)
}

func (h *AdminHandler) Raffles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	raffles, err := h.raffleRepo.FindAll(ctx)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := make([]adminRaffleItem, 0, len(raffles))
	for _, raffle := range raffles {
		organizerName := ""
		if user, err := h.userRepo.FindByID(ctx, raffle.OrganizerID); err == nil {
			organizerName = user.Name
		}

		soldTickets, _ := h.ticketRepo.CountByRaffleAndStatus(ctx, raffle.ID, model.TicketStatusPaid)
		paidTickets, _ := h.ticketRepo.CountByRaffleAndStatus(ctx, raffle.ID, model.TicketStatusPaid)
		revenue := soldTickets * int64(raffle.TicketPrice)

		items = append(items, adminRaffleItem{
			ID:            raffle.ID,
			OrganizerID:   raffle.OrganizerID,
			OrganizerName: organizerName,
			Title:         raffle.Title,
			TicketPrice:   raffle.TicketPrice,
			MaxNumbers:    raffle.MaxNumbers,
			DrawDate:      raffle.DrawDate,
			ImageURL:      raffle.ImageURL,
			Status:        string(raffle.Status),
			WinnerNumber:  raffle.WinnerNumber,
			CreatedAt:     raffle.CreatedAt,
			SoldTickets:   soldTickets,
			PaidTickets:   paidTickets,
			Revenue:       revenue,
		})
	}

	writeJSON(w, http.StatusOK, items)
}

type adminRaffleItem struct {
	ID              primitive.ObjectID `json:"id"`
	OrganizerID     primitive.ObjectID `json:"organizerId"`
	OrganizerName   string             `json:"organizerName"`
	Title           string             `json:"title"`
	TicketPrice     int                `json:"ticketPrice"`
	MaxNumbers      int                `json:"maxNumbers"`
	DrawDate        time.Time          `json:"drawDate"`
	ImageURL        string             `json:"imageUrl,omitempty"`
	Status          string             `json:"status"`
	WinnerNumber    *int               `json:"winnerNumber,omitempty"`
	CreatedAt       time.Time          `json:"createdAt"`
	SoldTickets     int64              `json:"soldTickets"`
	PaidTickets     int64              `json:"paidTickets"`
	Revenue         int64              `json:"revenue"`
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

	raffleRevenue, _ := h.paymentRepo.SumPaidByType(ctx, model.PaymentTypeRaffle)

	activeNonTrialUsers, _ := h.userRepo.FindActiveNonTrialUsers(ctx)
	subscriptionRevenue := int64(len(activeNonTrialUsers) * model.SubscriptionPrice)

	writeJSON(w, http.StatusOK, adminStats{
		TotalUsers:       totalUsers,
		ActiveUsers:      activeUsers,
		TotalRaffles:     totalRaffles,
		ActiveRaffles:    activeRaffles,
		TotalPaidTickets: totalPaidTickets,
		TotalRevenue:     raffleRevenue + subscriptionRevenue,
		TrialUsers:       trialUsers,
		PastDueUsers:     pastDueUsers,
	})
}
