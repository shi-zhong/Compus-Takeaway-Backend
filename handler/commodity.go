package handler

import (
	"backend/model"
	"backend/model/dbop"
	"backend/utils"
	"backend/utils/code"
	"github.com/gin-gonic/gin"
)

type CommodityCreateModel struct {
	Intro   string  `json:"intro"`
	Name    string  `json:"name"`
	Picture string  `json:"picture"`
	Price   float64 `json:"price"` // 商品单价
	Tags    string  `json:"tags"`
}

type CommodityUpdateModel struct {
	ID      uint    `json:"id"`
	Intro   string  `json:"intro"`
	Name    string  `json:"name"`
	Picture string  `json:"picture"`
	Price   float64 `json:"price"` // 商品单价
	Tags    string  `json:"tags"`
}

func CommodityCreateHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	commodityCreateModel := &CommodityCreateModel{}
	if !utils.QuickBind(c, commodityCreateModel) {
		return
	}

	if commodityCreateModel.Price <= 0 {
		code.GinBadRequest(c)
		return
	}

	commodity, msgCode, _ := dbop.CommodityInfoCreate(model.Db.Self, &model.CommodityInfo{
		ShopID:  token.IDExtra,
		Name:    commodityCreateModel.Name,
		Price:   commodityCreateModel.Price,
		Intro:   commodityCreateModel.Intro,
		Status:  0,
		Picture: commodityCreateModel.Picture,
		Tags:    commodityCreateModel.Tags,
	})

    if msgCode.Code == code.InsertError || msgCode.Code == code.DBEmpty {
		code.GinServerError(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"id":      commodity.ID,
		"status":  commodity.Status,
		"name":    commodity.Name,
		"intro":   commodity.Intro,
		"picture": commodity.Picture,
		"price":   commodity.Price,
	})
}
func CommodityUpdateHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	commodityUpdateModel := &CommodityUpdateModel{}
	if !utils.QuickBind(c, commodityUpdateModel) {
		return
	}

	if commodityUpdateModel.Price <= 0 {
		code.GinBadRequest(c)
		return
	}

	_, msgCode, _ := dbop.CommodityInfoUpdate(
		model.Db.Self,
		&model.CommodityInfo{
			ShopID: token.IDExtra,
			ID:     commodityUpdateModel.ID,
		},
		&model.CommodityInfo{
			Name:    commodityUpdateModel.Name,
			Price:   commodityUpdateModel.Price,
			Intro:   commodityUpdateModel.Intro,
			Picture: commodityUpdateModel.Picture,
			Tags:    commodityUpdateModel.Tags,
		})

	if msgCode.Code == code.UpdateError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinBadRequest(c)
		return
	}

	code.GinOKEmpty(c)
}

func CommodityDeleteHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	pathIDModel := &PathIDModel{}
	if !utils.QuickBindPath(c, pathIDModel) {
		return
	}

	_, msgCode, _ := dbop.CommodityInfoUpdate(
		model.Db.Self,
		&model.CommodityInfo{
			ID:     pathIDModel.ID,
			ShopID: token.IDExtra,
		},
		&model.CommodityInfo{
			Status: model.CommodityStatusDeleted,
		},
	)

	if msgCode.Code == code.UpdateError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinBadRequest(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"status": model.CommodityStatusDeleted,
	})
	return
}
func CommodityDetailHandler(c *gin.Context) {
	pathIDModel := &PathIDModel{}
	if !utils.QuickBindPath(c, pathIDModel) {
		return
	}

	commodity, msgCode, _ := dbop.CommodityInfoCheck(&model.CommodityInfo{
		ID: pathIDModel.ID,
        Status: 0,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinBadRequest(c)
		return
	}

	code.GinOKPayloadAny(c, commodity[0])
}

func CommodityShopListHandler(c *gin.Context) {

	pathIDModel := &PathIDModel{}
	if !utils.QuickBindPath(c, pathIDModel) {
		return
	}

	limit := c.Query("limit")
	page := c.Query("page")

	commoditys, msgCode, _ := dbop.CommodityInfoLimitPageCheck(&model.CommodityInfo{
		ShopID: pathIDModel.ID,
	}, limit, page)

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"list":  commoditys,
		"count": len(commoditys),
	})

}
