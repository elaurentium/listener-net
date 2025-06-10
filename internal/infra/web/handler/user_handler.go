package handler

import (
	"net/http"

	"github.com/elaurentium/listener-net/internal/domain/service"
	"github.com/gin-gonic/gin"
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
	IP          string `json:"ip" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Dispositive string `json:"dispositive" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.UserService.Register(c.Request.Context(), req.IP, req.Name, req.Dispositive)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID, exist := c.Get("user_id")
	
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.UserService.GetByID(c.Request.Context(), userID.(uuid.UUID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserDispositive(c *gin.Context) {
	userMac, exist := c.Get("user_dispositive")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.UserService.GetByDispositive(c.Request.Context(), userMac.(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
