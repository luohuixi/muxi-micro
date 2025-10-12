package main

import (
	"errors"
	"github.com/muxi-Infra/muxi-micro/pkg/errs"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/zapx"
)

var DBNOData = errs.NewErr("db_no_data", "db has no data")

var UserNotFound = errs.NewErr("user_not_found", "User not found")

func SearchDB(id int) error {
	return DBNOData.WithCause(errors.New("row is 0"))
}

func SearchUser(id int) error {
	err := SearchDB(id)
	if errors.Is(err, DBNOData) {
		return UserNotFound.WithCause(err).WithMeta(map[string]interface{}{
			"user_id": id,
		})
	}
	return nil
}

func main() {
	l := zapx.NewDefaultZapLogger()
	id := 1
	A(id, l)
	l.Info("查询出错id:", logger.Field{"id": id})
}

func A(id int, l logger.Logger) {
	err := SearchUser(id)
	l.Error("查询数据库出错:" + err.Error())
}
