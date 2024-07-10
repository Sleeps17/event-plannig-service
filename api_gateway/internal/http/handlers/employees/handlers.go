package employees

import (
	employeeservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/employees"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	logger            *slog.Logger
	employeesProvider employeeservice.Service
}

func New(logger *slog.Logger, employeesProvider employeeservice.Service) *Handler {
	return &Handler{logger: logger, employeesProvider: employeesProvider}
}

func (h *Handler) Register(router *gin.Engine) {
	employee := router.Group("/employees")

	employee.GET("/", h.GetAllEmployees)
	employee.GET("/:id", h.GetEmployeeByID)
	employee.PUT("/:id", h.UpdateEmployee)
	employee.DELETE("/:id", h.DeleteEmployeeByID)
}
