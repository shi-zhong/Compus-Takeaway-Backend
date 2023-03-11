package handler

import (
	"backend/utils"
	"backend/utils/code"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"strings"
	"time"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func UploadImg(c *gin.Context) {

	f, err := c.FormFile("imgfile")
	if err != nil {
		code.GinBadRequest(c)
		return
	} else {
		fileExt := strings.ToLower(path.Ext(f.Filename))
		if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
			code.GinBadRequest(c)
			return
		}
		fileName := utils.MD5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := "./static/"
		isExist := PathExists(fildDir)
		if !isExist {
			mkerr := os.Mkdir(fildDir, os.ModePerm)
            if mkerr!= nil {
                code.GinServerError(c)
                return
            }
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		saveErr := c.SaveUploadedFile(f, filepath)
        if saveErr != nil {
            code.GinServerError(c)
            return
        }

        code.GinOKPayload(c, &gin.H{
            "url": "http://localhost:8000" + filepath[1:],
        })
	}
}
