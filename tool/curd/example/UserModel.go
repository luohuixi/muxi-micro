package model

import (
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/sql"
	"time"
)

var _ UserModels = (*ExtraUserExec)(nil)

type UserModels interface {
	UserModel
	// 可以在这里添加额外的方法接口
}

type ExtraUserExec struct {
	*UserExec
}

func NewUserModels(DBdsn, redisAddr, redisPassword string, number int, ttlForCache, ttlForSet time.Duration, l logger.Logger) (UserModels, error) {
	db, err := sql.ConnectDB(DBdsn, User{})
	if err != nil {
		return nil, err
	}
	cache := sql.ConnectCache(redisAddr, redisPassword, number, ttlForCache, ttlForSet)

	instance := NewUserModel(db, cache, l)

	return &ExtraUserExec{
		instance,
	}, nil
}

// 可以在这里添加额外的方法实现
// func (e *ExtraUserExec) ...