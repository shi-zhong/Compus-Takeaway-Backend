package handler

import (
	"backend/model"
	"backend/model/dbop"
	"backend/utils"
	"backend/utils/code"
	"github.com/gin-gonic/gin"
)

type AddressModel struct {
	Receiver string `json:"receiver"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type AddressUpdateModel struct {
	ID       uint   `json:"id"`
	Receiver string `json:"receiver"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

func AddressAddHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	addressToAdd := &AddressModel{}
	if !utils.QuickBind(c, addressToAdd) {
		return
	}

	_, msgCode, _ := dbop.CustomerAddressCreate(model.Db.Self, &model.CustomerAddress{
		UserID:   token.ID,
		Address:  addressToAdd.Address,
		Phone:    addressToAdd.Phone,
		Receiver: addressToAdd.Receiver,
	})

	if msgCode.Code == code.InsertError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.InvalidPhone {
		code.GinBadRequest(c)
		return
	}

	code.GinOKEmpty(c)
}
func AddressUpdateHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	addressUpdateModel := &AddressUpdateModel{}
	if !utils.QuickBind(c, addressUpdateModel) {
		return
	}

	_, msgCode2, _ := dbop.CustomerAddressUpdate(
		model.Db.Self,
		&model.CustomerAddress{
			ID:     addressUpdateModel.ID,
			UserID: token.ID,
		},
		&model.CustomerAddress{
			Address:  addressUpdateModel.Address,
			Phone:    addressUpdateModel.Phone,
			Receiver: addressUpdateModel.Receiver,
		},
	)

	if msgCode2.Code == code.UpdateError {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}

func AddressCheckHander(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	addresses, msgCode, _ := dbop.CustomerAddressCheck(&model.CustomerAddress{
		UserID: token.ID,
	})

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	}

	var ades = make([]*AddressUpdateModel, len(addresses))

	for index, address := range addresses {
		ades[index] = &AddressUpdateModel{
			ID:       address.ID,
			Address:  address.Address,
			Phone:    address.Phone,
			Receiver: address.Receiver,
		}
	}

	code.GinOKPayload(c, &gin.H{
		"address": ades,
		"count":   len(ades),
	})
}

func AddressDeleteHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	addressPathModel := &PathIDModel{}
	utils.QuickBindPath(c, addressPathModel)

	msgCode, _ := dbop.CustomerAddressDrop(model.Db.Self, &model.CustomerAddress{
		ID:     addressPathModel.ID,
		UserID: token.ID,
	})

	if msgCode.Code == code.DropError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinOKEmpty(c)
		return
	}

	code.GinOKEmpty(c)
}
