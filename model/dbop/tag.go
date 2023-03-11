package dbop

import (
    "backend/model"
    "backend/utils/code"
    "gorm.io/gorm"
)

// TagCreate InsertError
func TagCreate(tx *gorm.DB, tag *model.Tag) (*model.Tag, *code.MsgCode, error) {

    // 数据库存储
    result := tx.Create(tag)

    // 插入有问题
    if result.Error != nil {
        return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
    }

    // 插入成功
    return tag, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}


func TagDrop(tx *gorm.DB, condition *model.Tag) (*code.MsgCode, error) {

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


func TagCheck(condition *model.Tag) ([]*model.Tag, *code.MsgCode, error) {

    var searchTag []*model.Tag

    // 条件由外部决定
    result := model.Db.Self.Where(condition).
        Not(&model.Tag{}).
        Find(&searchTag)

    if result.Error != nil {
        return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
    }

    // 找不到用户
    if result.RowsAffected == 0 {
        return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
    }

    return searchTag, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}


// TagUpdate  tag 就是要保存的数据
func TagUpdate(tx *gorm.DB, condition, tag *model.Tag) (*model.Tag, *code.MsgCode, error) {
    result := tx.Model(condition).Updates(tag)

    if result.Error != nil {
        return nil, &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
    }

    if result.RowsAffected == 0 {
        return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
    }

    return nil, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
