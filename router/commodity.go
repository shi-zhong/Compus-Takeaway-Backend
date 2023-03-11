package router

import (
	"backend/handler"
    "backend/middleware"
    "github.com/gin-gonic/gin"
)

func setCommodityRouter(router *gin.RouterGroup) {

    // 店铺内商品列表
    router.GET("/commodity/all/:ID", handler.CommodityShopListHandler)

	group := router.Group("/commodity")
    group.Use(middleware.ShopOnly())
    {
		group.POST("/create", handler.CommodityCreateHandler)
		group.POST("/update", handler.CommodityUpdateHandler)
		// 上架下架
		group.PUT("/del/:ID", handler.CommodityDeleteHandler)
		group.GET("/detail/:ID", handler.CommodityDetailHandler)
	}
}
