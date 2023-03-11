package router

import (
	"backend/handler"
	"github.com/gin-gonic/gin"
)

func setUploadRouter(router *gin.RouterGroup) {
	router.POST("/upload", handler.UploadImg)
}
