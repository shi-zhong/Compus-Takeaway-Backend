package dbop

import (
    "backend/model"
    "backend/utils/code"
    "gorm.io/gorm"
)

// PhysicalAddressCreate InsertError
func PhysicalAddressCreate(tx *gorm.DB, physicalAddress *model.PhysicalAddress) (*model.PhysicalAddress, *code.MsgCode, error) {

	/*

	   数据校验

	*/

	// 重复校验
	_, msgCode, err := PhysicalAddressCheck(&model.PhysicalAddress{})

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
	result := tx.Create(physicalAddress)

	// 插入有问题
	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	// 插入成功
	return physicalAddress, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}


func PhysicalAddressCheck(condition *model.PhysicalAddress) (*model.PhysicalAddress, *code.MsgCode, error) {

    var searchPhysicalAddress *model.PhysicalAddress

    // 条件由外部决定
    result := model.Db.Self.Where(condition).Find(&searchPhysicalAddress)

    if result.Error != nil {
        return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
    }

    // 找不到用户
    if result.RowsAffected == 0 {
        return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
    }

    return searchPhysicalAddress, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
