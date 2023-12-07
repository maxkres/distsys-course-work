package handler

import (
	"github.com/gin-gonic/gin"
	"store/config"
	"store/internal/client"
	"store/internal/service"
)

type Handler struct {
	s   *service.Service
	c   *client.Client
	Cfg *config.Config
}

func New(s *service.Service, cfg *config.Config) (*Handler, error) {
	client, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	return &Handler{
		s,
		client,
		cfg,
	}, nil
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	router.POST("put", h.Store)
	router.GET("get", h.Get)

	return router
}
