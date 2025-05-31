package handler

import (
	"encoding/json"
	"net/http"

	"github.com/elaurentium/listener-net/internal/domain/service"
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


func (h *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	var req UserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	
}
