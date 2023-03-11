package handler

import (
	"backend/model"
	"backend/model/dbop"
	"backend/utils"
	"backend/utils/code"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ShopInfoGet(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	queryID := c.Query("id")

	var searchID uint
	if queryID == "" {
		searchID = token.IDExtra
	} else {
		// string to uint
		queryID2, err := strconv.Atoi(queryID)
		if err != nil {
			code.GinServerError(c)
			return
		}
		searchID = uint(queryID2)
	}

	// search
	shopInfo, msgCode, _ := dbop.ShopInfoCheck(&model.Shop{
		ID: searchID,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinUserNotExist(c)
		return
	}

	// search Phone
	shopkeeper, msgCode2, _ := dbop.UserCheck(&model.User{
		ID: shopInfo.ShopKeeperID,
	})

	if msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.DBEmpty {
		code.GinUserNotExist(c)
		return
	}
	// search Phone
	shopAddress, msgCode3, _ := dbop.PhysicalAddressCheck(&model.PhysicalAddress{
		ID: shopInfo.AddressID,
	})

	if msgCode3.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode3.Code == code.DBEmpty {
		code.GinUserNotExist(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"id":     shopInfo.ID,
		"name":   shopInfo.ShopName,
		"intro":  shopInfo.ShopIntro,
		"avatar": shopInfo.ShopAvatar,
		"phone":  shopkeeper.Phone,
		"address": &gin.H{
			"building": shopAddress.BuildingID,
			"floor":    shopAddress.BuildingFloor,
			"number":   shopAddress.BuildingNumber,
		},
		"star":    4.5,
		"monthly": 600,
	})
}

type updateShopInfoModel struct {
	ShopName         string  `json:"name"`
	ShopIntro        string  `json:"intro"`
	ShopAvatar       string  `json:"avatar"`
	ShopStartDeliver float64 `json:"start_deliver"`
}

func ShopInfoPost(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	updateModel := &updateShopInfoModel{}
	if err := c.BindJSON(updateModel); err != nil {
		code.GinBadRequest(c)
		return
	}

	// search
	_, msgCode, _ := dbop.ShopInfoUpdate(
		model.Db.Self,
		&model.Shop{
			ID: token.IDExtra,
		},
		&model.Shop{
			ShopName:     updateModel.ShopName,
			ShopIntro:    updateModel.ShopIntro,
			ShopAvatar:   updateModel.ShopAvatar,
			StartDeliver: updateModel.ShopStartDeliver,
		})

	if msgCode.Code == code.UpdateError {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}

func ShopStatus(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	shopInfo, msgCode, _ := dbop.ShopInfoCheck(&model.Shop{
		ID: token.IDExtra,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinUserNotExist(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"status": shopInfo.CanBeSearched,
	})
}

type udpateShopStatus struct {
	Status uint `json:"status"`
}

func ShopStatusChange(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	updateModel := &udpateShopStatus{}
	if err := c.BindJSON(updateModel); err != nil {
		code.GinBadRequest(c)
		return
	}

	// search
	_, msgCode, _ := dbop.ShopInfoUpdate(
		model.Db.Self,
		&model.Shop{
			ID: token.IDExtra,
		},
		&model.Shop{
			CanBeSearched: updateModel.Status,
		})

	if msgCode.Code == code.UpdateError {
		code.GinServerError(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
        "status": updateModel.Status,
    })
}
