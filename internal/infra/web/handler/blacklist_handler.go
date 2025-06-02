package handler

import "github.com/elaurentium/listener-net/internal/domain/service"


type BlacklistHandler struct {
	BlacklistService *service.BlacklistService
}

type BlacklistRequest struct {
	IP 			string    `json:"ip" binding:"required"`
	Name		string    `json:"name" binding:"required"`
	Dispositive	string    `json:"dispositive" binding:"required"`
}

