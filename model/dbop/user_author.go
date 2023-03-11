package dbop

import (
	"backend/model"
	"backend/utils/code"
	"gorm.io/gorm"
)

func UserCreate(tx *gorm.DB, user *model.User) (*model.User, *code.MsgCode, error) {
	// 重复校验
	_, msgCode, err := UserCheck(&model.User{
		Phone:    user.Phone,
		Identity: user.Identity,
	})

	if msgCode.Code == code.CheckError {
		return nil, msgCode, err
	}

	// 数据库存储
	result := tx.Create(user)

	// 插入有问题
	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "InsertError", Code: code.InsertError}, result.Error
	}

	// 插入成功
	return user, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

// UserCheck CheckError; DBEmpty; OK;
func UserCheck(user *model.User) (*model.User, *code.MsgCode, error) {

	searchUser := &model.User{}

	// 条件由外部决定
	result := model.Db.Self.Where(user).Find(searchUser)

	if result.Error != nil {
		return nil, &code.MsgCode{Msg: "CheckError", Code: code.CheckError}, result.Error
	}

	// 找不到用户
	if result.RowsAffected == 0 {
		return nil, &code.MsgCode{Msg: "DBEmpty", Code: code.DBEmpty}, nil
	}

	return searchUser, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}

// UserUpdate user 就是要保存的数据
func UserUpdate(tx *gorm.DB, user *model.User) (*code.MsgCode, error) {

	result := tx.Save(user)

	if result.Error != nil {
		return &code.MsgCode{Msg: "UpdateError", Code: code.UpdateError}, result.Error
	}

	return &code.MsgCode{Msg: "OK", Code: code.OK}, nil
}
