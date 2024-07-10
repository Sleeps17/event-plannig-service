package events

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	eventservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/events"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateEventRequest struct {
	Event *models.Event `json:"event"`
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var request CreateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidRequestMsg})
		return
	}

	resp, err := h.eventsProvider.Create(c, request.Event)
	if err != nil {
		if errors.Is(err, eventservice.ErrInternalServer) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
			return
		}

		if errors.Is(err, eventservice.ErrRoomIsOccupied) {
			c.JSON(http.StatusBadRequest, gin.H{"error": handlers.RoomIsOccupiedMsg})
			return
		}

		if errors.Is(err, eventservice.ErrSomeWorkersAreBusy) {
			c.JSON(http.StatusConflict, resp.BusyEmployees)
			return
		}
	}

	c.JSON(http.StatusOK, resp.ID)
}
