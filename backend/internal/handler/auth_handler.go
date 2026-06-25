package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/user/rifa-online/internal/middleware"
	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type authResponse struct {
	User         *userResponse `json:"user"`
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
}

type userResponse struct {
	ID                 string                  `json:"id"`
	Name               string                  `json:"name"`
	Email              string                  `json:"email"`
	Role               model.Role              `json:"role"`
	Phone              string                  `json:"phone,omitempty"`
	InfinitePayHandle  string                  `json:"infinitePayHandle,omitempty"`
	SubscriptionStatus model.SubscriptionStatus `json:"subscriptionStatus"`
	SubscriptionIsTrial bool                   `json:"subscriptionIsTrial"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.authService.Register(r.Context(), service.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyRegistered) {
			writeError(w, "email already registered", http.StatusConflict)
			return
		}
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, toAuthResponse(result))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.authService.Login(r.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeError(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, toAuthResponse(result))
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	user, err := h.authService.GetProfile(r.Context(), userID)
	if err != nil {
		writeError(w, "user not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(user))
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID := middleware.UserIDFromContext(r.Context())
	user, err := h.authService.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyRegistered) {
			writeError(w, "email already in use", http.StatusConflict)
			return
		}
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(user))
}

func toUserResponse(user *model.User) *userResponse {
	status := user.SubscriptionStatus
	if status == "" {
		status = model.SubscriptionStatusInactive
	}
	role := user.Role
	if role == "" {
		role = model.RoleUser
	}
	return &userResponse{
		ID:                 user.ID.Hex(),
		Name:               user.Name,
		Email:              user.Email,
		Role:                role,
		Phone:               user.Phone,
		InfinitePayHandle:   user.InfinitePayHandle,
		SubscriptionStatus:  status,
		SubscriptionIsTrial: user.SubscriptionIsTrial,
	}
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		writeError(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, toAuthResponse(result))
}

func toAuthResponse(result *service.AuthResult) *authResponse {
	return &authResponse{
		User:         toUserResponse(result.User),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
