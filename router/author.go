package router

import (
	"backend/handler"
	"github.com/gin-gonic/gin"
)

func setAuthorRouterWithoutAuthorize(router *gin.RouterGroup) {
	group := router.Group("/author")
	{
        group.POST("/login", handler.LoginHandler)

    }
}


func setAuthorRouterWithAuthorize(router *gin.RouterGroup) {
    group := router.Group("/author")
    {
        group.POST("/shopkeeper/login", handler.ShopKeeperLoginHandler)
        group.POST("/customer/login", handler.OhterToCustomerLogin)
    }
}
