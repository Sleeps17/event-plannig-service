package events

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	eventservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/events"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UpdateEventRequest struct {
	Event *models.Event `json:"event"`
}

func (h *Handler) UpdateRoom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidIdParamMsg})
		return
	}

	var request UpdateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidRequestMsg})
		return
	}

	request.Event.ID = uint64(id)
	resp, err := h.eventsProvider.Update(c, request.Event)
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

		if errors.Is(err, eventservice.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": handlers.EventNotFoundMsg})
			return
		}
	}

	c.JSON(http.StatusOK, resp.UpdatedEvent)
}
