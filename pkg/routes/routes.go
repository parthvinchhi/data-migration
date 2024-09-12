package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/data-migration/pkg/handlers"
)

func Routes(router *gin.Engine) {
	router.Static("/static", "./pkg/static")
	router.LoadHTMLGlob("pkg/templates/*")

	router.GET("/", handlers.LoadHtml)
	router.POST("/migrate", handlers.MigrateHandler)
}
