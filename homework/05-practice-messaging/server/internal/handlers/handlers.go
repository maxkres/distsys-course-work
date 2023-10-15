package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"server/internal/models"
	"sync"
)

var id int
var mu sync.RWMutex

type Handler struct {
	svc   ServiceInterface
	files []string
}

type ServiceInterface interface {
	StoreImage(ctx context.Context, url string) error
}

func New(svc ServiceInterface, delivery chan string) *Handler {
	h := &Handler{
		svc:   svc,
		files: make([]string, 0),
	}
	go func() {
		for m := range delivery {
			mu.Lock()
			h.files = append(h.files, m)
			mu.Unlock()
		}
	}()
	return h
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	api := router.Group("/api/v1.0/")
	api.POST("/images", h.StoreImage)
	api.GET("/images", h.GetImages)
	api.GET("/images/:id", h.GetImageDescription)

	return router
}

func (h *Handler) StoreImage(c *gin.Context) {
	var req models.Requset
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.svc.StoreImage(c, fmt.Sprintf("%v:%v", id, req.ImageUrl)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"image_id": fmt.Sprint(id),
	})
	id++
}

func (h *Handler) GetImages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"image_ids": h.files,
	})
}

func (h *Handler) GetImageDescription(c *gin.Context) {
	id := c.Param("id")
	for _, file := range h.files {
		if file == id {
			body, _ := os.ReadFile("data/" + id)
			c.JSON(http.StatusOK, gin.H{
				"description": string(body),
			})
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)

}
