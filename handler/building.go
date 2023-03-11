package handler

import (
	"backend/model"
	"backend/model/dbop"
	"backend/utils/code"
	"github.com/gin-gonic/gin"
)

func AddBuilding(c *gin.Context) {}

func RemoveBuilding(c *gin.Context) {}
func EditBuilding(c *gin.Context)   {}

func GetBuilding(c *gin.Context) {
	buildings, msgCode, _ := dbop.BuildingCheck(&model.Building{})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	}

	code.GinOKPayload(c, &gin.H{
		"buildings": buildings,
		"count":     len(buildings),
	})
}
