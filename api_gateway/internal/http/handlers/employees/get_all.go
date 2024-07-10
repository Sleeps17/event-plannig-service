package employees

import (
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetAllEmployees(c *gin.Context) {
	employees, err := h.employeesProvider.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
		return
	}

	c.JSON(http.StatusOK, employees)
}
