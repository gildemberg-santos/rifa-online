package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/user/rifa-online/internal/config"
	"github.com/user/rifa-online/internal/middleware"
	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/service"
)

const (
	refreshCookieName   = "refresh_token"
	refreshCookieMaxAge = 7 * 24 * 60 * 60 // 7 dias, igual ao TTL do refresh token
	refreshCookiePath   = "/api/v1/auth"
)

type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authService: authService, cfg: cfg}
}

// setRefreshCookie grava o refresh token num cookie HttpOnly (inacessível ao JS,
// protegendo contra XSS). Secure é ligado fora de desenvolvimento; SameSite=Lax
// mitiga CSRF.
func (h *AuthHandler) setRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     refreshCookiePath,
		MaxAge:   refreshCookieMaxAge,
		HttpOnly: true,
		Secure:   h.cfg.AppEnv != "development",
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *AuthHandler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cfg.AppEnv != "development",
		SameSite: http.SameSiteLaxMode,
	})
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
	User        *userResponse `json:"user"`
	AccessToken string        `json:"accessToken"`
	// O refresh token NÃO é mais devolvido no corpo: ele vive apenas no cookie HttpOnly.
}

type userResponse struct {
	ID                  string                   `json:"id"`
	Name                string                   `json:"name"`
	Email               string                   `json:"email"`
	Role                model.Role               `json:"role"`
	Phone               string                   `json:"phone,omitempty"`
	InfinitePayHandle   string                   `json:"infinitePayHandle,omitempty"`
	SubscriptionStatus  model.SubscriptionStatus `json:"subscriptionStatus"`
	SubscriptionIsTrial bool                     `json:"subscriptionIsTrial"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if 	err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"user":       toUserResponse(result.User),
		"emailSent":  true,
		"message":    "Codigo de verificacao enviado para o email",
	})
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
		if errors.Is(err, service.ErrEmailNotVerified) {
			writeError(w, "email not verified", http.StatusForbidden)
			return
		}
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.setRefreshCookie(w, result.RefreshToken)
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
		ID:                  user.ID.Hex(),
		Name:                user.Name,
		Email:               user.Email,
		Role:                role,
		Phone:               user.Phone,
		InfinitePayHandle:   user.InfinitePayHandle,
		SubscriptionStatus:  status,
		SubscriptionIsTrial: user.SubscriptionIsTrial,
	}
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Preferência: cookie HttpOnly. Fallback ao corpo apenas por compatibilidade.
	token := ""
	if c, err := r.Cookie(refreshCookieName); err == nil {
		token = c.Value
	}
	if token == "" {
		var req refreshRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		token = req.RefreshToken
	}
	if token == "" {
		writeError(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	result, err := h.authService.RefreshToken(r.Context(), token)
	if err != nil {
		h.clearRefreshCookie(w)
		writeError(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	h.setRefreshCookie(w, result.RefreshToken)
	writeJSON(w, http.StatusOK, toAuthResponse(result))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.clearRefreshCookie(w)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type verifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type resendCodeRequest struct {
	Email string `json:"email"`
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req verifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.authService.VerifyEmail(r.Context(), req.Email, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCode):
			writeError(w, "codigo invalido", http.StatusBadRequest)
		case errors.Is(err, service.ErrCodeExpired):
			writeError(w, "codigo expirado, solicite um novo", http.StatusBadRequest)
		default:
			writeError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	h.setRefreshCookie(w, result.RefreshToken)
	writeJSON(w, http.StatusOK, toAuthResponse(result))
}

func (h *AuthHandler) ResendCode(w http.ResponseWriter, r *http.Request) {
	var req resendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.authService.ResendCode(r.Context(), req.Email); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "codigo reenviado"})
}

func toAuthResponse(result *service.AuthResult) *authResponse {
	return &authResponse{
		User:        toUserResponse(result.User),
		AccessToken: result.AccessToken,
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
