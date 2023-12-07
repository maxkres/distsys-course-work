package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) Store(c *gin.Context) {
	token, err := c.Cookie("jwt")
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	username, err := h.c.CheckJWT(c, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var value struct {
		Value string `json:"value"`
	}
	if err := c.BindJSON(&value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	key := c.Query("key")

	err = h.s.Store(username, key, value.Value)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (h *Handler) Get(c *gin.Context) {
	token, err := c.Cookie("jwt")
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	username, err := h.c.CheckJWT(c, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	key := c.Query("key")

	val, err := h.s.Get(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !strings.Contains(val, username) {
		c.Status(http.StatusForbidden)
		return
	}

	fmt.Println(val, username)
	c.JSON(http.StatusOK, gin.H{
		"value": strings.Split(val, ":")[1],
	})
}
