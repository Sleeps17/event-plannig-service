package rooms

import (
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	roomservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/rooms"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) DeleteRoom(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.InvalidIdParamMsg})
		return
	}

	if err := h.roomsProvider.Delete(c, uint32(id)); err != nil {
		if errors.Is(err, roomservice.ErrInternalServer) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
			return
		}

		if errors.Is(err, roomservice.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": handlers.RoomNotFoundMsg})
			return
		}
	}

	c.Status(http.StatusOK)
}
