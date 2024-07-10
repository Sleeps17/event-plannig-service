package rooms

import (
	roomservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/rooms"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	logger        *slog.Logger
	roomsProvider roomservice.Service
}

func New(logger *slog.Logger, roomsProvider roomservice.Service) *Handler {
	return &Handler{logger: logger, roomsProvider: roomsProvider}
}

func (h *Handler) Register(router *gin.Engine) {
	rooms := router.Group("/rooms")

	rooms.POST("/", h.CreateRoom)
	rooms.GET("/:id", h.GetRoomByID)
	rooms.GET("/", h.GetAllRooms)
	rooms.PUT("/:id", h.UpdateRoom)
	rooms.DELETE("/:id", h.DeleteRoom)
}
