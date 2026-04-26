package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	token := csrf.Token(c.Request)

	c.Header("X-CSRF-Token", token)

	c.JSON(200, gin.H{
		"user": nil,
		"role": "guest",
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// реализация
}
