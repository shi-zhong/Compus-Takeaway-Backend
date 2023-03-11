package router

import (
    "backend/handler"
    "github.com/gin-gonic/gin"
)

func setBuildingRouter(router *gin.RouterGroup) {
    group := router.Group("/building")
    {
        group.GET("/list", handler.GetBuilding)

    }
}
