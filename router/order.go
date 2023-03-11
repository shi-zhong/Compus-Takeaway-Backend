package router

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

func setOrderRouter(router *gin.RouterGroup) {
    router.GET("/order/detail/:ID", handler.OrderDetailHandler)


    orderGroupM := router.Group("/order")
	orderGroupM.Use(middleware.ShopOnly())
	{
		// 商家发货
        orderGroupM.GET("/list/shop", handler.OrderMerchentListHandler)
        orderGroupM.POST("/shop/accept/:ID", handler.OrderShopAcceptHandler)
        orderGroupM.POST("/shop/cookfinish/:ID", handler.OrderCookFinishHandler)
    }

	orderGroupC := router.Group("/order")
	orderGroupC.Use(middleware.CustomerOnly())
	{
		// 订单列表获取
		orderGroupC.GET("/list/customer", handler.OrderCustomerListHandler)
		// 获取订单细节
        orderGroupC.POST("/create", handler.OrderCreate)
        orderGroupC.POST("/cancel/:ID", handler.OrderCustomerCancelHandler)
        orderGroupC.POST("/finish/:ID", handler.OrderCustomerFinishHandler)

	}
}
