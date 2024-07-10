package events

import (
	eventservice "github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/services/events"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	logger         *slog.Logger
	eventsProvider eventservice.Service
}

func New(logger *slog.Logger, eventsProvider eventservice.Service) *Handler {
	return &Handler{logger: logger, eventsProvider: eventsProvider}
}

func (h *Handler) Register(router *gin.Engine) {
	events := router.Group("/events")

	events.POST("/", h.CreateEvent)
	events.GET("/", h.GetAllRooms)
	events.GET("/:id", h.GetEventByID)
	events.PUT("/:id", h.UpdateRoom)
	events.DELETE("/:id", h.DeleteRoom)
}
