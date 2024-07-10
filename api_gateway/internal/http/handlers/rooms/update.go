package rooms

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/domain/models"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	roomservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/rooms"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UpdateRoomRequest struct {
	Room *models.Room `json:"room"`
}

func (h *Handler) UpdateRoom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidIdParamMsg})
		return
	}

	var request UpdateRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidRequestMsg})
		return
	}

	request.Room.ID = uint32(id)
	updatedRoom, err := h.roomsProvider.Update(c, request.Room)
	if err != nil {
		if errors.Is(err, roomservice.ErrInternalServer) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
			return
		}

		if errors.Is(err, roomservice.ErrRoomAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": handlers.RoomAlreadyExists})
			return
		}

		if errors.Is(err, roomservice.ErrRoomNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": handlers.RoomNotFoundMsg})
			return
		}
	}

	c.JSON(http.StatusOK, updatedRoom)
}
