package handler

import (
	"backend/model"
	"backend/model/dbop"
	"backend/utils"
	"backend/utils/code"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type orderUnionModel struct {
	Order     *model.Order
	Commodity []*model.CommodityInfo
	Shop      *model.Shop
}

type basicCommodityModel struct {
	ID    uint `json:"id"`
	Count uint `json:"count"`
}

type orderCreateModel struct {
	Shop      uint                   `json:"shop"`
	Commodity []*basicCommodityModel `json:"commodity"`
	Address   uint                   `json:"address"`
}

func buildCommodityString(commodityModel []*basicCommodityModel) string {
	str := ""
	for _, commodity := range commodityModel {
		str += strconv.Itoa(int(commodity.ID)) + "/" + strconv.Itoa(int(commodity.Count)) + ";"
	}
	return str[:len(str)-1]
}

func OrderCreate(c *gin.Context) {
	token := utils.GetTokenInfo(c)
	if token.Identity != model.IdentityCustomer {
		code.GinUnAuthorized(c)
		return
	}

	order := &orderCreateModel{}
	if !utils.QuickBind(c, order) {
		code.GinBadRequest(c)
		return
	}

	shop, msgCode6, _ := dbop.ShopInfoCheck(&model.Shop{
		ID: order.Shop,
	})

	if msgCode6.Code != code.OK {
		code.GinBadRequest(c)
		return
	}

	if shop.CanBeSearched == model.ShopIsClose {
		code.GinOKPayload(c, &gin.H{
			"status": 0,
		})
	}

	shouldPay := 0.0

	// 校验所有商品
	for _, commodity := range order.Commodity {
		coo, msgCode, _ := dbop.CommodityInfoCheck(&model.CommodityInfo{
			ID:     commodity.ID,
			ShopID: order.Shop,
		})
		if msgCode.Code == code.CheckError {
			code.GinServerError(c)
			return
		} else if msgCode.Code == code.DBEmpty {
			code.GinMissingShareBill(c)
			return
		}
		shouldPay += coo[0].Price * float64(commodity.Count)
	}

	// 校验地址
	address, msgCode2, _ := dbop.CustomerAddressCheck(&model.CustomerAddress{
		ID:     order.Address,
		UserID: token.ID,
	})

	if msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.DBEmpty {
		code.GinMissingshop(c)
		return
	}

	// pay
	tx := model.Db.Self.Begin()

	checkuser, msgCode6, _ := dbop.UserCheck(&model.User{
		ID: token.ID,
	})

	if msgCode6.Code == code.CheckError {
		tx.Rollback()
		code.GinServerError(c)
		return
	}

	checkuser.LastUsedAddress = order.Address

	msgCode4, _ := dbop.UserUpdate(tx, checkuser)

	if msgCode4.Code == code.UpdateError {
		tx.Rollback()
		code.GinServerError(c)
		return
	}

	_, msgCode3, _ := dbop.OrderCreate(tx, &model.Order{
		ID:          utils.OrderIDGenerate(),
		CustomerID:  token.ID,
		ShopID:      order.Shop,
		RiderID:     6,
		CommodityID: buildCommodityString(order.Commodity),

		Receiver: address[0].Receiver,
		Phone:    address[0].Phone,
		Address:  address[0].Address,

		CreateAt: time.Now(),
		Status:   model.OrderCreate,
	})

	if msgCode3.Code == code.InsertError || msgCode3.Code == code.DBEmpty {
		tx.Rollback()
		code.GinServerError(c)
		return
	}

	tx.Commit()
	code.GinOKEmpty(c)
}

/*
返回店铺基本信息，
商品基本信息，
*/

func orderUnionCheck(order *model.Order) (*orderUnionModel, *code.MsgCode) {

	/// shop
	shop, msgCode, _ := dbop.ShopInfoCheck(&model.Shop{
		ID: order.ShopID,
	})

	if msgCode.Code != code.OK {
		return nil, msgCode
	}

	// commodity
	commoditys := strings.Split(order.CommodityID, ";")

	// 获取 商品id 和 数量

	ids := make([]uint, len(commoditys))
	counts := make([]uint, len(commoditys))

	for index, commodity := range commoditys {
		split := strings.Split(commodity, "/")
		id, err1 := strconv.Atoi(split[0])
		count, err2 := strconv.Atoi(split[1])

		if err1 != nil || err2 != nil {
			return nil, &code.MsgCode{Code: code.CheckError, Msg: "CheckError"}
		}

		ids[index] = uint(id)
		counts[index] = uint(count)
	}

	commodityInfos, msgCode2, _ := dbop.CommodityInfoCheckIn(&ids)
	if msgCode2.Code != code.OK {
		return nil, msgCode2
	}

	// address

	return &orderUnionModel{
		Order:     order,
		Shop:      shop,
		Commodity: commodityInfos,
	}, &code.MsgCode{Code: code.OK, Msg: "OK"}
}

func OrderCustomerListHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	limit := c.Query("limit")
	page := c.Query("page")

	orders, msgCode, _ := dbop.OrderLimitPageCheck(&model.Order{CustomerID: token.ID}, limit, page)

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	}
	var orderDetails []*orderUnionModel = make([]*orderUnionModel, len(orders))

	for index, order := range orders {
		orderUnion, msgCode2 := orderUnionCheck(order)
		if msgCode2.Code != code.OK {
			return
		}
		orderDetails[index] = orderUnion
	}

	code.GinOKPayload(c, &gin.H{
		"list":  orderDetails,
		"count": len(orderDetails),
	})
}

func OrderMerchentListHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	limit := c.Query("limit")
	page := c.Query("page")

	orders, msgCode, _ := dbop.OrderLimitPageCheck(&model.Order{ShopID: token.IDExtra}, limit, page)

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	}
	var orderDetails []*orderUnionModel = make([]*orderUnionModel, len(orders))

	for index, order := range orders {
		orderUnion, msgCode2 := orderUnionCheck(order)
		if msgCode2.Code != code.OK {
			return
		}
		orderDetails[index] = orderUnion
	}

	code.GinOKPayload(c, &gin.H{
		"list":  orderDetails,
		"count": len(orderDetails),
	})
}

func OrderDetailHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	path := &PathStringIDModel{}

	if !utils.QuickBindPath(c, path) {
		code.GinBadRequest(c)
		return
	}

	condition := &model.Order{
		ID: path.ID,
	}

	if token.Identity == model.IdentityCustomer {
		condition.CustomerID = token.ID
	} else if token.Identity == model.IdentityShopKeeper {
		condition.ShopID = token.IDExtra
	}

	orders, msgCode, _ := dbop.OrderCheck(condition)

	if msgCode.Code == code.CheckError {
		code.GinServerError(c)
		return
	}

	orderUnion, msgCode2 := orderUnionCheck(orders[0])
	if msgCode2.Code != code.OK {
		return
	}

	code.GinOKPayload(c, &gin.H{
		"order": orderUnion,
	})
}

func OrderShopAcceptHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	pathStringIDModel := &PathStringIDModel{}
	err := utils.QuickBindPath(c, pathStringIDModel)
	if !err {
		return
	}

	order, msgCode2, _ := dbop.OrderCheck(&model.Order{
		ShopID: token.IDExtra,
		ID:     pathStringIDModel.ID,
	})

	if msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.DBEmpty {
		code.GinMissingOrder(c)
		return
	}

	if order[0].Status != model.OrderCreate {
		code.GinOKEmpty(c)
		return
	}

	_, msgCode, _ := dbop.OrderUpdate(model.Db.Self, &model.Order{
		ShopID: token.IDExtra,
		ID:     pathStringIDModel.ID,
	}, &model.Order{
		AcceptAt: time.Now(),
		Status:   model.OrderAccept,
	})

	if msgCode.Code == code.UpdateError || msgCode.Code == code.DBEmpty {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}
func OrderCookFinishHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	pathStringIDModel := &PathStringIDModel{}
	err := utils.QuickBindPath(c, pathStringIDModel)
	if !err {
		return
	}

	order, msgCode2, _ := dbop.OrderCheck(&model.Order{
		ShopID: token.IDExtra,
		ID:     pathStringIDModel.ID,
	})

	if msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.DBEmpty {
		code.GinMissingOrder(c)
		return
	}

	if order[0].Status != model.OrderAccept {
		code.GinOKEmpty(c)
		return
	}

	_, msgCode, _ := dbop.OrderUpdate(model.Db.Self, &model.Order{
		ShopID: token.IDExtra,
		ID:     pathStringIDModel.ID,
	}, &model.Order{
		CookFinishAt: time.Now(),
		Status:       model.OrderCookFinish,
	})

	if msgCode.Code == code.UpdateError || msgCode.Code == code.DBEmpty {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}

func OrderCustomerCancelHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	pathStringIDModel := &PathStringIDModel{}
	err := utils.QuickBindPath(c, pathStringIDModel)
	if !err {
		return
	}

	order, msgCode2, _ := dbop.OrderCheck(&model.Order{
		CustomerID: token.ID,
		ID:         pathStringIDModel.ID,
	})

	if msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.DBEmpty {
		code.GinMissingOrder(c)
		return
	}

	if order[0].Status != model.OrderCreate && order[0].Status != model.OrderAccept {
		code.GinOKEmpty(c)
		return
	}

	_, msgCode, _ := dbop.OrderUpdate(model.Db.Self, &model.Order{
		CustomerID: token.ID,
		ID:         pathStringIDModel.ID,
	}, &model.Order{
		FinishAt: time.Now(),
		Status:   model.OrderCancel,
	})

	if msgCode.Code == code.UpdateError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinMissingOrder(c)
		return
	}

	if msgCode.Code == code.UpdateError || msgCode.Code == code.DBEmpty {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}

func OrderCustomerFinishHandler(c *gin.Context) {
	token := utils.GetTokenInfo(c)

	pathStringIDModel := &PathStringIDModel{}
	err := utils.QuickBindPath(c, pathStringIDModel)
	if !err {
		return
	}

	order, msgCode2, _ := dbop.OrderCheck(&model.Order{
		CustomerID: token.ID,
		ID:         pathStringIDModel.ID,
	})

	if msgCode2.Code == code.CheckError {
		code.GinServerError(c)
		return
	} else if msgCode2.Code == code.DBEmpty {
		code.GinMissingOrder(c)
		return
	}

	if order[0].Status != model.OrderDeliverFinish {
		code.GinOKEmpty(c)
		return
	}

	_, msgCode, _ := dbop.OrderUpdate(model.Db.Self, &model.Order{
		CustomerID: token.ID,
		ID:         pathStringIDModel.ID,
	}, &model.Order{
		FinishAt: time.Now(),
		Status:   model.OrderFinish,
	})

	if msgCode.Code == code.UpdateError {
		code.GinServerError(c)
		return
	} else if msgCode.Code == code.DBEmpty {
		code.GinMissingOrder(c)
		return
	}

	if msgCode.Code == code.UpdateError || msgCode.Code == code.DBEmpty {
		code.GinServerError(c)
		return
	}

	code.GinOKEmpty(c)
}
