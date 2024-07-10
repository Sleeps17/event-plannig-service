package auth

import (
	authservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/auth"
	employeeservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/employees"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	logger            *slog.Logger
	authProvider      authservice.Service
	employeesProvider employeeservice.Service
}

func New(logger *slog.Logger, authProvider authservice.Service, employeesProvider employeeservice.Service) *Handler {
	return &Handler{
		logger:            logger,
		authProvider:      authProvider,
		employeesProvider: employeesProvider,
	}
}

func (h *Handler) Register(router *gin.Engine) {
	auth := router.Group("/auth")

	auth.POST("/register", h.RegisterUser)
	auth.POST("/login", h.LoginUser)
}
