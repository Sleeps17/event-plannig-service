package events

import (
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GetAllRequest struct {
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
}

func (h *Handler) GetAllRooms(c *gin.Context) {
	var request GetAllRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(400, gin.H{"error": handlers.InvalidRequestMsg})
		return
	}

	employees, err := h.eventsProvider.GetAll(c, request.StartDate, request.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
		return
	}

	c.JSON(http.StatusOK, employees)
}
