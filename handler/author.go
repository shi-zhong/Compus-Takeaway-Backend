package handler

import (
	"backend/model"
	"backend/model/dbop"
	"backend/utils"
	"backend/utils/code"
	"github.com/gin-gonic/gin"
)

/*
顾客登录  自动创建账号
店家登录  需要提前注册(即申请)
管理员登录 直接登录，无需创建和申请
骑手登录 需注册
*/

type LoginModel struct {
	OpenID   string `json:"open_id"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	NickName string `json:"nickname"`
	Identity uint   `json:"identity"`
}

func LoginHandler(c *gin.Context) {
	// 绑定参数
	login := &LoginModel{}
	if err := c.BindJSON(login); err != nil {
		code.GinBadRequest(c)
	}

	// 检索数据库中有关信息
	checkUser, msgCode, _ := dbop.UserCheck(&model.User{
		OpenID:   login.OpenID,
		Identity: model.IdentityCustomer,
	})

	if msgCode.Code != code.OK && msgCode.Code != code.DBEmpty {
		// 不知名错误
		code.GinEmptyMsgCode(c, msgCode)
		return
	} else if msgCode.Code == code.DBEmpty {
		/*
		   新建用户
		   根据用户身份处理
		   顾客 -> 自动注册
		   店家，骑手，经理 -> 返回错误
		*/
		if login.Identity == model.IdentityCustomer {
			newUser := &model.User{
				OpenID:   login.OpenID,
				Phone:    login.Phone,
				Avatar:   login.Avatar,
				NickName: login.NickName,
				Identity: model.IdentityCustomer,
			}
			tx := model.Db.Self.Begin()
			user, msgCode2, _ := dbop.UserCreate(tx, newUser)

			if msgCode2.Code != code.OK {
				code.GinServerError(c)
			} else {
				checkUser = user
			}
			tx.Commit()
		} else {
			code.GinBadRequest(c)
			return
		}
	}

	/**
	  根据身份处理不同信息

	  店家额外返回店铺号
	*/

	var idextra uint = 0

	if login.Identity == model.IdentityShopKeeper {
		searchShop, msgCode3, _ := dbop.ShopInfoCheck(&model.Shop{
			ShopKeeperID: checkUser.ID,
		})
		if msgCode3.Code != code.OK {
			code.GinBadRequest(c)
			return
		}
		idextra = searchShop.ID
	}

	// 验证成功， 开始生成token
	token, err2 := utils.TokenEecode(&utils.TokenPayload{
		Identity: checkUser.Identity,
		ID:       checkUser.ID,
		Phone:    checkUser.Phone,
		IDExtra:  idextra,
	})

	if err2 != nil {
		code.GinServerError(c)
		return
	}

	// 发送token
	code.GinOKPayload(c, &gin.H{
		"token": token,
		"login": checkUser,
	})

}

func OhterToCustomerLogin(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	// 查询店家信息
	checkUser, msgCode, _ := dbop.UserCheck(&model.User{
		Phone:    token.Phone,
		Identity: model.IdentityCustomer,
	})

    if msgCode.Code !=code.OK {
        code.GinUnAuthorized(c)
        return
    }

	// 验证成功， 开始生成token
	token2, err2 := utils.TokenEecode(&utils.TokenPayload{
		Identity: checkUser.Identity,
		ID:       checkUser.ID,
		Phone:    checkUser.Phone,
		IDExtra:  0,
	})

	if err2 != nil {
		code.GinServerError(c)
		return
	}

	// 发送token
	code.GinOKPayload(c, &gin.H{
		"token":      token2,
		"user": checkUser,
	})
}

func checkShopKeeperIfAllowed(phone string) bool {
	return "18396148343" == phone
}

func checkShopKeeperAddressID(phone string) uint {
	return 1
}

func ShopKeeperLoginHandler(c *gin.Context) {

	token := utils.GetTokenInfo(c)

	// 查询店家信息
	checkUser, msgCode, _ := dbop.UserCheck(&model.User{
		Phone:    token.Phone,
		Identity: model.IdentityShopKeeper,
	})

	shop := &model.Shop{}

	if msgCode.Code != code.OK && msgCode.Code != code.DBEmpty {
		// 不知名错误
		code.GinEmptyMsgCode(c, msgCode)
		return
	} else if msgCode.Code == code.DBEmpty {
		// 未查询到 准备新建用户

		// 查询资格
		if !checkShopKeeperIfAllowed(token.Phone) {
			code.GinUnAuthorized(c)
			return
		}

		// 查询基本信息
		checkCus, msgCode2, _ := dbop.UserCheck(&model.User{
			ID: token.ID,
		})

		if msgCode2.Code != code.OK {
			// 不知名错误
			code.GinEmptyMsgCode(c, msgCode)
			return
		}
		// 新建用户
		newUser := &model.User{
			OpenID:   checkCus.OpenID,
			Phone:    checkCus.Phone,
			Avatar:   checkCus.Avatar,
			NickName: checkCus.NickName,
			Identity: model.IdentityShopKeeper,
		}
		tx := model.Db.Self.Begin()
		user, msgCode3, _ := dbop.UserCreate(tx, newUser)

		shop2, msgCode4, _ := dbop.ShopInfoCreate(tx, &model.Shop{
			ShopKeeperID:  user.ID,
			ShopName:      "Default Shop Name",
			ShopIntro:     "Default Shop Introduce",
			ShopAvatar:    "https://st-gdx.dancf.com/gaodingx/0/uxms/design/20210812-184716-154c.png",
			AddressID:     checkShopKeeperAddressID(token.Phone),
			CanBeSearched: 0,
		})

		if msgCode3.Code != code.OK || msgCode4.Code != code.OK {
			tx.Rollback()
			code.GinServerError(c)
		} else {
			checkUser = user
		}

		shop = shop2

		tx.Commit()
	}

	shop3, msgCode4, _ := dbop.ShopInfoCheck(&model.Shop{
		ShopKeeperID: checkUser.ID,
	})
	if msgCode4.Code != code.OK {
		code.GinServerError(c)
	} else {
		shop = shop3
	}

	// 验证成功， 开始生成token
	token2, err2 := utils.TokenEecode(&utils.TokenPayload{
		Identity: checkUser.Identity,
		ID:       checkUser.ID,
		Phone:    checkUser.Phone,
		IDExtra:  shop.ID,
	})

	if err2 != nil {
		code.GinServerError(c)
		return
	}

	// 发送token
	code.GinOKPayload(c, &gin.H{
		"token":      token2,
		"shopkeeper": checkUser,
		"shop_id":    shop.ID,
	})
}

/**
  登录 各种身份

  @todo 获取店铺信息
*/

type loginByWxModel struct {
	Code     string
	Phone    string
	NickName string
	Avatar   string
	Identity uint
}

func LoginByWxHandler(c *gin.Context) {
	// 通过 code 请求 session_key 和 open_id
	openId := "openid"
	// 解析数据获得手机号，昵称和头像信息
	user := &loginByWxModel{}
	// 参数绑定失败
	if err := c.BindJSON(user); err != nil {
		code.GinBadRequest(c)
		return
	}

	// 检索数据库中有关信息
	checkUser, msgCode, _ := dbop.UserCheck(&model.User{
		OpenID:   openId,
		Identity: user.Identity,
	})

	if msgCode.Code != code.OK && msgCode.Code != code.DBEmpty {
		// 不知名错误
		code.GinEmptyMsgCode(c, msgCode)
		return
	} else if msgCode.Code == code.DBEmpty {
		// 新建用户
		newUser := &model.User{}
		tx := model.Db.Self.Begin()
		user, msgCode2, _ := dbop.UserCreate(tx, newUser)

		if msgCode2.Code != code.OK {
			code.GinServerError(c)
		} else {
			checkUser = user
		}
	}

	// 验证成功， 开始生成token
	token, err2 := utils.TokenEecode(&utils.TokenPayload{
		Identity: checkUser.Identity,
		ID:       checkUser.ID,
		Phone:    checkUser.Phone,
	})

	if err2 != nil {
		code.GinServerError(c)
		return
	}

	// 发送token
	code.GinOKPayload(c, &gin.H{
		"token": token,
	})
}
