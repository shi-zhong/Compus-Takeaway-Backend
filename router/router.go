package router

import (
    "backend/middleware"
    "github.com/gin-gonic/gin"
)

var Router *gin.Engine

func Init() {
	Router = gin.Default()

    Router.Use(middleware.Cors())

    Router.Static("/static", "./static")

    // 不需要token
    donotNeedAuthorize := Router.Group("/api/v1")
    {

        setAuthorRouterWithoutAuthorize(donotNeedAuthorize)
        setSearchRouter(donotNeedAuthorize)
        setBuildingRouter(donotNeedAuthorize)
    }

	needAuthorize := Router.Group("/api/v1")
    needAuthorize.Use(middleware.TokenAuthorize())
	{
        setAuthorRouterWithAuthorize(needAuthorize)
        setShopRouter(needAuthorize)
        setUserInfoRouter(needAuthorize)
        setCommodityRouter(needAuthorize)
        setOrderRouter(needAuthorize)
        setUploadRouter(donotNeedAuthorize)
	}
}
