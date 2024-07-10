package auth

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	authservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterRequest struct {
	Login    string           `json:"login"`
	Password string           `json:"password"`
	Employee *models.Employee `json:"employee"`
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidRequestMsg})
	}

	if err := h.authProvider.Register(
		c,
		request.Login,
		request.Password,
	); err != nil {
		if errors.Is(err, authservice.ErrInternalServer) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
			return
		}

		if errors.Is(err, authservice.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidRequestMsg})
			return
		}

		if errors.Is(err, authservice.ErrUserAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": handlers.UserAlreadyExistsMsg})
			return
		}
	}

	id, err := h.employeesProvider.Create(c, request.Employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	request.Employee.ID = id
	c.JSON(http.StatusOK, request.Employee)
}
