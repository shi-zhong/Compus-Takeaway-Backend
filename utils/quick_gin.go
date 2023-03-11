package utils

import (
	"backend/utils/code"
	"github.com/gin-gonic/gin"
)

// GetTokenInfo ID Identity Phone
func GetTokenInfo(c *gin.Context) *TokenPayload {
	ID, exist := c.Get("ID")
	uintID, _ := ID.(uint)

	Identity, exist2 := c.Get("Identity")
	uintIdentity, _ := Identity.(uint)

	Phone, exist3 := c.Get("Phone")
	stringPhone, _ := Phone.(string)

	IDExtra, exist4 := c.Get("IDExtra")
	uintExtra, _ := IDExtra.(uint)

	if !exist || !exist2 || !exist3 || !exist4 {
		code.GinUnAuthorized(c)
		c.Abort()
		return nil
	}
	return &TokenPayload{
		ID:       uintID,
		Identity: uintIdentity,
		Phone:    stringPhone,
		IDExtra:  uintExtra,
	}
}

func QuickBind(c *gin.Context, structPointer any) bool {
	if err := c.BindJSON(structPointer); err != nil {
		code.GinBadRequest(c)
		c.Abort()
		return false
	}
	return true
}

func QuickBindPath(c *gin.Context, structPointer any) bool {
	if err := c.ShouldBindUri(structPointer); err != nil {
		code.GinBadRequest(c)
		c.Abort()
		return false
	}
	return true
}
