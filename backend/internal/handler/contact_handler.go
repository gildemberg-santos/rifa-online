package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/user/rifa-online/internal/model"
	"github.com/user/rifa-online/internal/repository"
)

type ContactHandler struct {
	contactRepo *repository.ContactRepo
}

func NewContactHandler(contactRepo *repository.ContactRepo) *ContactHandler {
	return &ContactHandler{contactRepo: contactRepo}
}

type contactRequest struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Message string `json:"message"`
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req contactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Contact = strings.TrimSpace(req.Contact)
	req.Message = strings.TrimSpace(req.Message)

	if req.Name == "" || len(req.Name) > 150 {
		writeError(w, "nome é obrigatório (máx. 150 caracteres)", http.StatusBadRequest)
		return
	}
	if len(req.Contact) > 150 {
		writeError(w, "contato muito longo (máx. 150 caracteres)", http.StatusBadRequest)
		return
	}
	if l := len(req.Message); l < 10 || l > 2000 {
		writeError(w, "mensagem deve ter entre 10 e 2000 caracteres", http.StatusBadRequest)
		return
	}

	msg := &model.ContactMessage{
		Name:    req.Name,
		Contact: req.Contact,
		Message: req.Message,
		IP:      r.RemoteAddr,
	}
	if err := h.contactRepo.Insert(r.Context(), msg); err != nil {
		writeError(w, "não foi possível registrar a mensagem", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (h *ContactHandler) List(w http.ResponseWriter, r *http.Request) {
	messages, err := h.contactRepo.ListRecent(r.Context(), 200)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, messages)
}
