package main

import (
	"errors"
	"fmt"
	"github.com/muxi-Infra/muxi-micro/pkg/errs"
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
	err := SearchUser(1)
	fmt.Println(err)

	// res
	// [user_not_found] User not found map[user_id:1] => [db_no_data] db has no data => row is 0
}
