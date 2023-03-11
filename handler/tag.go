package handler

import (
	"backend/model"
	"backend/model/dbop"
	"backend/utils"
	"backend/utils/code"
	"github.com/gin-gonic/gin"
)

type tagModel struct {
	Tag string
	ID  uint
}

func TagAddHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	tag := &tagModel{}
	if !utils.QuickBind(c, tag) {
		return
	}

	if tag.Tag != "" {
		_, msgCode, _ := dbop.TagCreate(model.Db.Self, &model.Tag{
			Tag:    tag.Tag,
			Belong: token.IDExtra,
		})
		if msgCode.Code == code.InsertError {
			code.GinServerError(c)
			return
		}

		code.GinOKEmpty(c)
		return
	}
	code.GinBadRequest(c)
}

func TagDelHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	tag := &tagModel{}
	if !utils.QuickBind(c, tag) {
		return
	}

	msgCode, _ := dbop.TagDrop(model.Db.Self, &model.Tag{
		ID:     tag.ID,
		Belong: token.IDExtra,
	})

	if msgCode.Code == code.DropError {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
	return
}

func TagUpdHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	tag := &tagModel{}
	if !utils.QuickBind(c, tag) {
		return
	}

	if tag.Tag != "" {
		_, msgCode, _ := dbop.TagUpdate(model.Db.Self,
			&model.Tag{
				ID:     tag.ID,
				Belong: token.IDExtra,
			},
			&model.Tag{
				Tag: tag.Tag,
			})
		if msgCode.Code == code.UpdateError {
			code.GinServerError(c)
			return
		}

		code.GinOKEmpty(c)
		return
	}
	code.GinBadRequest(c)
}

func TagAllHandler(c *gin.Context) {

	path := &PathIDModel{}
	if !utils.QuickBindPath(c, path) {
		return
	}

	tags, _, _ := dbop.TagCheck(&model.Tag{
		Belong: path.ID,
	})

	code.GinOKPayload(c, &gin.H{
		"tags": tags,
	})
}
