package utils

import (
	"backend/utils/code"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type TokenPayload struct {
	Identity uint
	ID       uint
	Phone    string
	IDExtra  uint
}
type TokenPayloadClaims struct {
	TokenPayload
	jwt.RegisteredClaims // 注册当前结构体为 Claims
}

func TokenEecode(payload *TokenPayload) (string, error) {
	// 创建秘钥
	key := []byte(GlobalConfig.SecretKey.Private)

	// 创建Token结构体
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenPayloadClaims{
		TokenPayload: *payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(GlobalConfig.GinConfig.TokenExpires) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})
	// 调用加密方法，发挥Token字符串
	signedString, err := token.SignedString(key)
	return signedString, err
}

func TokenDecode(signedString string) (*TokenPayloadClaims, *code.MsgCode, error) {
	// 根据Token字符串解析成Claims结构体
	token, err := jwt.ParseWithClaims(signedString, &TokenPayloadClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GlobalConfig.SecretKey.Private), nil
	})

	// 简要错误处理
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 || ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, &code.MsgCode{
					Msg: "TokenInvalid", Code: code.TokenInvalid,
				}, err
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, &code.MsgCode{Msg: "TokenExpired", Code: code.TokenExpired}, err
			} else {
				return nil, &code.MsgCode{Msg: "ServerError", Code: code.ServerError}, err
			}
		}
	}

	if claims, ok := token.Claims.(*TokenPayloadClaims); ok && token.Valid {
		return claims, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
	}
	return nil, &code.MsgCode{Msg: "ServerError", Code: code.ServerError}, err
}
