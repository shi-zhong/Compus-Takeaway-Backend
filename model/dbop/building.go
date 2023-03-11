package dbop

import (
	"backend/model"
	"backend/utils/code"
	"gorm.io/gorm"
	"strconv"
)

func buildingFloorStringToUintArray(floor string) []uint {
	return nil
}

func buildingFloorUintArrayToString(floor []uint) string {
	return ""
}

// BuildingCreate InsertError
func BuildingCreate(tx *gorm.DB, building *model.Building) (*model.Building, *code.MsgCode, error) {

	/*

	   数据校验

	*/

	// 重复校验
	_, msgCode, err := BuildingCheck(&model.Building{})

	if msgCode.Code == code.CheckError {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, err
	}

	if msgCode.Code == code.OK {
		return nil, &code.MsgCode{}, nil
	}

	if msgCode.Code == code.DBEmpty {
		return nil, &code.MsgCode{}, nil
	}

	/*

	   额外数据处理

	*/

	// 数据库存储
	result := tx.Create(building)

	// 插入有问题
	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	// 插入成功
	return building, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func BuildingDrop(tx *gorm.DB, condition *model.Building) (*code.MsgCode, error) {

	result := tx.Delete(condition)

	if result.Error != nil {
		return &code.MsgCode{Msg: "DropError", Code: code.DropError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func BuildingCheck(condition *model.Building) ([]*model.Building, *code.MsgCode, error) {

	var searchBuilding []*model.Building

	// 条件由外部决定
	result := model.Db.Self.Where(condition).
		Not(&model.Building{}).
		Find(&searchBuilding)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return searchBuilding, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func BuildingLimitPageCheck(condition *model.Building, limit, page string) ([]*model.Building, *code.MsgCode, error) {

	var searchBuilding []*model.Building

	limitInt, _ := strconv.Atoi(limit)
	pageInt, _ := strconv.Atoi(page)

	if limitInt == 0 || pageInt == 0 {
		return BuildingCheck(condition)
	}

	// 条件由外部决定
	result := model.Db.Self.
		Where(condition).
		Not(&model.Building{}).
		Limit(limitInt).
		Offset(limitInt * pageInt).
		Find(searchBuilding)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return searchBuilding, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

// BuildingUpdate  building 就是要保存的数据
func BuildingUpdate(tx *gorm.DB, condition, building *model.Building) (*model.Building, *code.MsgCode, error) {
	result := tx.Model(condition).Updates(building)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return nil, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
