package employees

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	employeeservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/employees"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UpdateEmployeeRequest struct {
	Employee *models.Employee `json:"employee"`
}

type UpdateEmployeeResponse struct {
	UpdatedEmployee *models.Employee `json:"updated_employee"`
}

func (h *Handler) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidIdParamMsg})
		return
	}

	var request UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, handlers.InvalidRequestMsg)
		return
	}

	request.Employee.ID = uint64(id)

	updatedEmployee, err := h.employeesProvider.Update(c, request.Employee)
	if err != nil {
		if errors.Is(err, employeeservice.ErrInternalServer) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
			return
		}

		if errors.Is(err, employeeservice.ErrEmployeeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": handlers.EmployeeNotFoundMsg})
			return
		}

		if errors.Is(err, employeeservice.ErrEmployeeAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": handlers.EmployeeAlreadyExists})
			return
		}
	}

	c.JSON(http.StatusOK, UpdateEmployeeResponse{UpdatedEmployee: updatedEmployee})
}
