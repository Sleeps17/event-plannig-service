package auth

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	authservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	AppID    uint32 `json:"app_id"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) LoginUser(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authProvider.Login(c, request.Login, request.Password, request.AppID)
	if err != nil {
		if errors.Is(err, authservice.ErrInternalServer) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
			return
		}

		if errors.Is(err, authservice.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidCredentialsMsg})
			return
		}
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}
