package handler

import (
	"backend/model/dbop"
	"backend/utils"
	"backend/utils/code"
	"github.com/gin-gonic/gin"
)



func SearchHandler(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")

	search := &dbop.SearchModel{}
	if !utils.QuickBind(c, search) {
		code.GinBadRequest(c)
		return
	}

	shops := dbop.SearchCommoditiesLimitPage(limit, page, search)

	code.GinOKPayload(c, &gin.H{
		"shops": shops,
	})
}
