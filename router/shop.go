package router

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

func setShopRouter(router *gin.RouterGroup) {

	group := router.Group("/shop")
	{
		// 用户信息查询
		shop := group.Group("/info")
		{
			// 顾客基本信息获取
			shop.GET("/get", handler.ShopInfoGet)
		}
		tag := group.Group("/tag")
		{
            tag.GET("/all/:ID", handler.TagAllHandler)
		}

	}
	groupShopKeeper := router.Group("/shop")
	groupShopKeeper.Use(middleware.ShopOnly())
	{
        shop := groupShopKeeper.Group("/info")
		{
			shop.POST("/update", handler.ShopInfoPost)
		}
        tag := groupShopKeeper.Group("/tag")
        {
            tag.POST("/add",handler.TagAddHandler)
            tag.POST("/upd",handler.TagUpdHandler)
            tag.PUT("/del",handler.TagDelHandler)
        }
        groupShopKeeper.GET("/status", handler.ShopStatus)
        groupShopKeeper.POST("/status", handler.ShopStatusChange)
    }

}
