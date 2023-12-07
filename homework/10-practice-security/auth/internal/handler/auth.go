package handler

import (
	"auth/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) SignUp(c *gin.Context) {

	var user model.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := h.s.Check(user); err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	token, err := h.s.SignUp(user)
	if err != nil {
		//if errors.Is(err, service.ErrUserAlreadyExists) {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"error": err.Error(),
		//	})
		//	return
		//}
		//
		//logger.Error("/users/auth/sing-up", zap.Error(fmt.Errorf("service sing up failed: %w", err)))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("jwt", token, 0, "/", "", false, true)
	c.Status(http.StatusOK)
}

func (h *Handler) LogIn(c *gin.Context) {

	var user model.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	token, err := h.s.SingIn(user)
	if err != nil {
		//if errors.Is(err, service.ErrUserDoesNotExists) || errors.Is(err, service.ErrIncorrectPassword) {
		//
		//	c.JSON(http.StatusForbidden, gin.H{
		//		"error": err.Error(),
		//	})
		//	return
		//}
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("jwt", token, 0, "/", "", false, true)
	c.Status(http.StatusOK)
}

func (h *Handler) VerifyToken(c *gin.Context) {
	token, err := c.Cookie("jwt")
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	username, err := h.s.Verify(token)
	if err != nil {
		//if errors.Is(err, service.ErrTokenExpired) {
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		//		"error": err.Error(),
		//	})
		//	return
		//}
		//if strings.Contains(err.Error(), jwt.ErrSignatureInvalid.Error()) {
		//	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		//		"error": fmt.Errorf("wrong signature").Error(),
		//	})
		//	return
		//}
		//
		//logger.Error("service verify failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("verify failed: %w", err).Error(),
		})
		return
	}

	if err := h.s.Check(model.User{
		username,
		"",
	}); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Hello, %v", username))
}
