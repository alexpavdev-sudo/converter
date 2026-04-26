package controllers

import (
	"converter/app"
	"converter/services/formater"
	"github.com/gin-gonic/gin"
)

func GetFormats(c *gin.Context) {
	formatService := formater.NewFormatService()
	app.OK(c, formatService.GetFormats())
}
