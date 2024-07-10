package httpserver

import (
	"context"
	"fmt"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/config"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/handlers"
	"github.com/Sleeps17/events-planning-service-backend/api_gateway/internal/http/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
}

func New(cfg *config.Config, handlers ...handlers.Handler) *Server {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.AuthMiddleware(cfg.Secret))

	for _, handler := range handlers {
		handler.Register(r)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{server: server}
}

func (r *Server) Run() error {
	return r.server.ListenAndServe()
}

func (r *Server) ShutDown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.server.Shutdown(ctx)
}
