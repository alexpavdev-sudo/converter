package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"net/http"
)

func Ping(c *gin.Context) {
	session := sessions.Default(c)

	var count int
	v := session.Get("count")
	if v == nil {
		count = 1
	} else {
		count = v.(int)
		count++
	}

	session.Set("count", count)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session save failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"count":     count,
		"csrfToken": csrf.Token(c.Request),
	})
}
