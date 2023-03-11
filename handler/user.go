package handler

import (
	"backend/model"
	"backend/model/dbop"
	"backend/utils"
	"backend/utils/code"
	"github.com/gin-gonic/gin"
)

type UpdateCustomerModel struct {
	Avatar    *string `json:"avatar"`
	Introduce *string `json:"introduce"`
	Nickname  *string `json:"nickname"`
}

type UpdateshopModel struct {
	Address    string  `json:"address"`
	Avatar     *string `json:"avatar"`
	Introduce  *string `json:"introduce"`
	Nickname   string  `json:"nickname"`
	ShopAvatar *string `json:"shop_avatar"`
	ShopIntro  *string `json:"shop_intro"`
}

func GetSelfBasicInfo(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	// search
	user, msgCode, _ := dbop.UserCheck(&model.User{
		ID: token.ID,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinUserNotExist(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"user_id":  user.ID,
		"nickname": user.NickName,
		"avatar":   user.Avatar,
		"phone":    utils.HideMobile(user.Phone),
        "last_used_address": user.LastUsedAddress,	
	})
}

func ToBeShopKeeper(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	_, msgCode, _ := dbop.UserCheck(&model.User{
		Phone:    token.Phone,
		Identity: model.IdentityShopKeeper,
	})

	checkUser, msgCode2, _ := dbop.UserCheck(&model.User{
		Phone:    token.Phone,
		Identity: model.IdentityCustomer,
	})

	if msgCode.Code == code.CheckError || msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty && msgCode2.Code == code.OK {

		tx := model.Db.Self.Begin()

		insert, msgCode3, _ := dbop.UserCreate(tx, &model.User{
			Identity: model.IdentityShopKeeper,
			Phone:    checkUser.Phone,
			NickName: checkUser.NickName,
			Avatar:   checkUser.Avatar,
		})
		if msgCode3.Code == code.InsertError {
			code.GinServerError(c)
			tx.Rollback()
			return
		}

		_, msgCode4, _ := dbop.ShopInfoCreate(tx, &model.Shop{
			ShopKeeperID: insert.ID,
			ShopName:     "Shop Name",
			ShopIntro:    "",
			ShopAvatar:   "",
			AddressID:    0,
		})
		if msgCode4.Code == code.InsertError {
			code.GinServerError(c)
			tx.Rollback()
			return
		}
		tx.Commit()
	}

	code.GinOKEmpty(c)
}
