package employees

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	employeeservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/employees"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) GetEmployeeByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidIdParamMsg})
		return
	}

	employee, err := h.employeesProvider.GetByID(c, uint64(id))
	if err != nil {
		if errors.Is(err, employeeservice.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": handlers.EmployeeNotFoundMsg})
			return
		}

		if errors.Is(err, employeeservice.ErrInternalServer) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
			return
		}
	}

	c.JSON(http.StatusOK, employee)
}
