package handler

import (
	"auth/config"
	"auth/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	s   *service.Service
	Cfg *config.Config
}

func New(s *service.Service, cfg *config.Config) (*Handler, error) {
	return &Handler{s, cfg}, nil
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	router.POST("signup", h.SignUp)
	router.POST("login", h.LogIn)
	router.GET("whoami", h.VerifyToken)

	return router
}
