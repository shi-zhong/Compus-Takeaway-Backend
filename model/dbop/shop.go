package dbop

import (
    "backend/model"
    "backend/utils/code"
    "gorm.io/gorm"
)

func ShopInfoCreate(tx *gorm.DB, shop *model.Shop) (*model.Shop, *code.MsgCode, error) {
    // 重复校验
    _, msgCode, err := ShopInfoCheck(&model.Shop{
        ID: shop.ID,
        })

    //  数据库出错
    if msgCode.Code == code.CheckError {
        return nil, msgCode, err
    }

    if msgCode.Code == code.OK {
        // 转到更新操作
        return nil, nil, nil
    }

    // 数据库存储
    result := tx.Create(shop)

    // 插入有问题
    if result.Error != nil {
        return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
    }

    // 插入成功
    return shop, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

func ShopInfoCheck(shop *model.Shop) (*model.Shop, *code.MsgCode, error) {
    searchshop := &model.Shop{}

    // 条件由外部决定
    result := model.Db.Self.Where(shop).Find(searchshop)

    if result.Error != nil {
        return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
    }

    // 找不到商家
    if result.RowsAffected == 0 {
        return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
    }

    return searchshop, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}


func ShopInfoUpdate(tx *gorm.DB, condition, shop *model.Shop) (*model.Shop, *code.MsgCode, error) {
    result := tx.Model(condition).Updates(shop)

    if result.Error != nil {
        return nil, &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
    }

    if result.RowsAffected == 0 {
        return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
    }

    return nil, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
