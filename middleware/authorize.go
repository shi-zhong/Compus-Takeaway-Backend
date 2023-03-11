package middleware

import (
	"backend/model"
	"backend/utils"
	"backend/utils/code"
	_ "fmt"
	"github.com/gin-gonic/gin"
)

func TokenAuthorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// prehandle
		token := c.GetHeader("Token")

		tokenPayloadClaims, msgCode, _ := utils.TokenDecode(token)

		if msgCode.Code == code.TokenInvalid || msgCode.Code == code.ServerError {
			code.GinEmptyMsgCode(c, msgCode)
			c.Abort()
			return
		}

		if msgCode.Code == code.TokenExpired {
			// 后期改成更新模式
			code.GinEmptyMsgCode(c, msgCode)
			c.Abort()
			return
		}

		c.Set("ID", tokenPayloadClaims.TokenPayload.ID)
		c.Set("Identity", tokenPayloadClaims.TokenPayload.Identity)
		c.Set("Phone", tokenPayloadClaims.TokenPayload.Phone)
		c.Set("IDExtra", tokenPayloadClaims.TokenPayload.IDExtra)

		c.Next()

		// afterhandle

	}
}

func CustomerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := utils.GetTokenInfo(c)

		if token.Identity != model.IdentityCustomer {
			code.GinUnAuthorized(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

func ShopOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := utils.GetTokenInfo(c)

		if token.Identity != model.IdentityShopKeeper {
			code.GinUnAuthorized(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
