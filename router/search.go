package router

import (
	"backend/handler"
	"github.com/gin-gonic/gin"
)

func setSearchRouter(router *gin.RouterGroup) {
	router.POST("/search", handler.SearchHandler)
}
