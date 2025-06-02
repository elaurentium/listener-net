package handler

import (
	"encoding/json"
	"net/http"

	"github.com/elaurentium/listener-net/internal/domain/service"
	"github.com/google/uuid"
)


type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}


type UserRequest struct {
	IP 		string    `json:"ip" binding:"required"`
	Name	string    `json:"name" binding:"required"`
}


func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userRaw := r.Context().Value("user_ip")

	if userRaw == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
	}

	userID, ok := userRaw.(uuid.UUID)
	
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	user, err := h.UserService.GetByID(r.Context(), userID)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
