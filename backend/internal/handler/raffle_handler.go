package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/user/rifa-online/internal/middleware"
	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
	"github.com/user/rifa-online/internal/service"
)

type RaffleHandler struct {
	raffleService *service.RaffleService
	userRepo      *repository.UserRepo
}

func NewRaffleHandler(raffleService *service.RaffleService, userRepo *repository.UserRepo) *RaffleHandler {
	return &RaffleHandler{raffleService: raffleService, userRepo: userRepo}
}

type createRaffleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TicketPrice int    `json:"ticketPrice"`
	MaxNumbers  int    `json:"maxNumbers"`
	DrawDate    string `json:"drawDate"`
	ImageURL    string `json:"imageUrl,omitempty"`
}

func (h *RaffleHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req createRaffleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	drawDate, err := time.Parse(time.RFC3339, req.DrawDate)
	if err != nil {
		writeError(w, "invalid drawDate format (use RFC3339)", http.StatusBadRequest)
		return
	}

	raffle, err := h.raffleService.Create(r.Context(), service.CreateRaffleInput{
		OrganizerID: userID,
		Title:       req.Title,
		Description: req.Description,
		TicketPrice: req.TicketPrice,
		MaxNumbers:  req.MaxNumbers,
		DrawDate:    drawDate,
		ImageURL:    req.ImageURL,
	})
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, raffle)
}

func (h *RaffleHandler) List(w http.ResponseWriter, r *http.Request) {
	raffles, err := h.raffleService.ListActive(r.Context())
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, raffles)
}

func (h *RaffleHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writeError(w, "invalid raffle id", http.StatusBadRequest)
		return
	}

	userID := middleware.UserIDFromContext(r.Context())

	detail, err := h.raffleService.GetDetail(r.Context(), oid, userID)
	if err != nil {
		if errors.Is(err, service.ErrRaffleNotFound) {
			writeError(w, "raffle not found", http.StatusNotFound)
			return
		}
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

func (h *RaffleHandler) isAdmin(ctx context.Context, userID string) bool {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false
	}
	user, err := h.userRepo.FindByID(ctx, oid)
	if err != nil {
		return false
	}
	return user.Role == model.RoleAdmin
}

func (h *RaffleHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writeError(w, "invalid raffle id", http.StatusBadRequest)
		return
	}

	organizerOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req createRaffleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	drawDate, err := time.Parse(time.RFC3339, req.DrawDate)
	if err != nil {
		writeError(w, "invalid drawDate format", http.StatusBadRequest)
		return
	}

	isAdmin := h.isAdmin(r.Context(), userID)

	raffle, err := h.raffleService.Update(r.Context(), oid, organizerOID, service.CreateRaffleInput{
		Title:       req.Title,
		Description: req.Description,
		TicketPrice: req.TicketPrice,
		MaxNumbers:  req.MaxNumbers,
		DrawDate:    drawDate,
		ImageURL:    req.ImageURL,
	}, isAdmin)
	if err != nil {
		if errors.Is(err, service.ErrRaffleNotFound) {
			writeError(w, "raffle not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrNotRaffleOwner) {
			writeError(w, "not the raffle owner", http.StatusForbidden)
			return
		}
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, raffle)
}

func (h *RaffleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writeError(w, "invalid raffle id", http.StatusBadRequest)
		return
	}

	organizerOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	isAdmin := h.isAdmin(r.Context(), userID)

	if err := h.raffleService.Delete(r.Context(), oid, organizerOID, isAdmin); err != nil {
		if errors.Is(err, service.ErrRaffleNotFound) {
			writeError(w, "raffle not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrNotRaffleOwner) {
			writeError(w, "not the raffle owner", http.StatusForbidden)
			return
		}
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *RaffleHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writeError(w, "invalid raffle id", http.StatusBadRequest)
		return
	}

	organizerOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	isAdmin := h.isAdmin(r.Context(), userID)

	if err := h.raffleService.Cancel(r.Context(), oid, organizerOID, isAdmin); err != nil {
		if errors.Is(err, service.ErrRaffleNotFound) {
			writeError(w, "raffle not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrNotRaffleOwner) {
			writeError(w, "not the raffle owner", http.StatusForbidden)
			return
		}
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (h *RaffleHandler) MyRaffles(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	raffles, err := h.raffleService.GetMyRaffles(r.Context(), oid)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, raffles)
}

func (h *RaffleHandler) Stats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writeError(w, "invalid raffle id", http.StatusBadRequest)
		return
	}

	organizerOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	isAdmin := h.isAdmin(r.Context(), userID)

	stats, err := h.raffleService.GetStats(r.Context(), oid, organizerOID, isAdmin)
	if err != nil {
		if errors.Is(err, service.ErrRaffleNotFound) {
			writeError(w, "raffle not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrNotRaffleOwner) {
			writeError(w, "not the raffle owner", http.StatusForbidden)
			return
		}
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (h *RaffleHandler) DashboardStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	stats, err := h.raffleService.GetDashboardStats(r.Context(), oid)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (h *RaffleHandler) Draw(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		writeError(w, "invalid raffle id", http.StatusBadRequest)
		return
	}

	organizerOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		writeError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	isAdmin := h.isAdmin(r.Context(), userID)

	result, err := h.raffleService.Draw(r.Context(), oid, organizerOID, isAdmin)
	if err != nil {
		if errors.Is(err, service.ErrRaffleNotFound) {
			writeError(w, "raffle not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrNotRaffleOwner) {
			writeError(w, "not the raffle owner", http.StatusForbidden)
			return
		}
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
