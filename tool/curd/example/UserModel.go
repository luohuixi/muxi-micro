package model

import (
	"github.com/muxi-Infra/muxi-micro/pkg/sql"
)

var _ UserModels = (*ExtraUserExec)(nil)

type UserModels interface {
	UserModel
	// 可以在这里添加额外的方法接口
}

type ExtraUserExec struct {
	*UserExec
}

func NewUserModels(DBdsn string) (UserModels, error) {
	db, err := sql.ConnectDB(DBdsn, User{})
	if err != nil {
		return nil, err
	}
	instance := NewUserModel(db)

	return &ExtraUserExec{
		instance,
	}, nil
}

// 可以在这里添加额外的方法实现
// func (e *ExtraUserExec) ...