package template

import (
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/zapx"
	"github.com/muxi-Infra/muxi-micro/pkg/sql"
	"time"
)

var _ UserModels = (*ExtraUserExec)(nil)

type UserModels interface {
	UserModel
	// ...
}

type ExtraUserExec struct {
	*UserExec
}

func NewUserModels(DBdsn, redisAddr, redisPassword string, number int, ttl, ttl2 time.Duration) (UserModels, error) {
	db, err := sql.ConnectDB(DBdsn, User{})
	if err != nil {
		return nil, err
	}
	cache := sql.ConnectCache(redisAddr, redisPassword, number, ttl, ttl2)
	l := zapx.NewDefaultZapLogger(logger.EnvTest, true, "./logs")
	instance := NewUserModel(db, cache, l)
	return &ExtraUserExec{
		instance,
	}, nil
}

// func (e *ExtraUserModel)
