package dbop

import (
	"backend/model"
	"backend/utils/code"
	"gorm.io/gorm"
	"strconv"
)

func CommodityInfoCreate(tx *gorm.DB, commodity *model.CommodityInfo) (*model.CommodityInfo, *code.MsgCode, error) {

	// 数据库存储
	result := tx.Create(commodity)

	// 插入有问题
	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	// 插入成功
	return commodity, &code.MsgCode{Msg: "OK", Code: code.OK}, nil

}

// CommodityInfoCheck CheckError; DBEmpty; OK;
func CommodityInfoCheck(commodity *model.CommodityInfo) ([]*model.CommodityInfo, *code.MsgCode, error) {

	var searchCommodityInfo []*model.CommodityInfo

	// 条件由外部决定
	result := model.Db.Self.Where(commodity).
		Not(&model.CommodityInfo{Status: model.CommodityStatusDeleted}).
		Find(&searchCommodityInfo)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return searchCommodityInfo, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func CommodityInfoCheckIn(ids *[]uint) ([]*model.CommodityInfo, *code.MsgCode, error) {

	var searchCommodityInfo []*model.CommodityInfo

	// 条件由外部决定
	result := model.Db.Self.Find(&searchCommodityInfo, ids)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return searchCommodityInfo, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func CommodityInfoLimitPageCheck(commodity *model.CommodityInfo, limit, page string) ([]*model.CommodityInfo, *code.MsgCode, error) {

	var searchCommodityInfo []*model.CommodityInfo

	limitInt, _ := strconv.Atoi(limit)
	pageInt, _ := strconv.Atoi(page)

	if limitInt == 0 || pageInt == 0 {
		return CommodityInfoCheck(commodity)
	}

	// 条件由外部决定
	result := model.Db.Self.
		Where(commodity).
		Not(&model.CommodityInfo{Status: model.CommodityStatusDeleted}).
		Limit(limitInt).
		Offset(limitInt * pageInt).
		Find(searchCommodityInfo)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return searchCommodityInfo, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func CommodityInfoUpdate(tx *gorm.DB, condition, commodity *model.CommodityInfo) (*model.CommodityInfo, *code.MsgCode, error) {
	result := tx.Model(condition).Updates(commodity)

	tx.Debug().Model(condition).Updates(commodity)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return commodity, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
