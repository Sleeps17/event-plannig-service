package rooms

import (
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetAllRooms(c *gin.Context) {
	rooms, err := h.roomsProvider.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": handlers.InternalErrorMsg})
		return
	}

	c.JSON(http.StatusOK, rooms)
}
