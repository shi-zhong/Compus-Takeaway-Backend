package router

import (
	"backend/handler"
    "backend/middleware"
    "github.com/gin-gonic/gin"
)

func setUserInfoRouter(router *gin.RouterGroup) {
	group := router.Group("/user")
	{
		// 用户信息查询
		userInfo := group.Group("/info")
		{
			// 顾客基本信息获取
			userInfo.GET("/customer/get", handler.GetSelfBasicInfo)
		}

		userAddress := group.Group("/address")
        userAddress.Use(middleware.CustomerOnly())
		{
            // 用户only
			userAddress.POST("/add", handler.AddressAddHandler)
            userAddress.GET("/list", handler.AddressCheckHander)
			userAddress.PUT("/update", handler.AddressUpdateHandler)
			userAddress.DELETE("/delete/:ID", handler.AddressDeleteHandler)
		}

	}
}
