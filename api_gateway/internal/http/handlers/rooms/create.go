package rooms

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	roomservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/rooms"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateRoomRequest struct {
	Room *models.Room `json:"room"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var request CreateRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InternalErrorMsg})
		return
	}

	id, err := h.roomsProvider.Create(c, request.Room)
	if err != nil {
		if errors.Is(err, roomservice.ErrRoomAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": handlers.RoomAlreadyExists})
			return
		}

		if errors.Is(err, roomservice.ErrInternalServer) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
			return
		}
	}

	request.Room.ID = id
	c.JSON(http.StatusOK, request.Room)
}
