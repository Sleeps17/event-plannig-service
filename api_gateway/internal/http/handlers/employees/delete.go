package employees

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	employeeservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/employees"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) DeleteEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidIdParamMsg})
		return
	}

	if err := h.employeesProvider.Delete(c, uint64(id)); err != nil {
		if errors.Is(err, employeeservice.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": handlers.EmployeeNotFoundMsg})
			return
		}

		if errors.Is(err, employeeservice.ErrInternalServer) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
			return
		}
	}

	c.Status(http.StatusOK)
}
