package events

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	eventservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/events"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) GetEventByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidIdParamMsg})
		return
	}

	event, err := h.eventsProvider.GetByID(c, uint64(id))
	if err != nil {
		if errors.Is(err, eventservice.ErrInternalServer) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
			return
		}

		if errors.Is(err, eventservice.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": handlers.EventNotFoundMsg})
			return
		}
	}

	c.JSON(http.StatusOK, event)
}
