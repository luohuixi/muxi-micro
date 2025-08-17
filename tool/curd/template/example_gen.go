package template

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/sql"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

const (
	cacheUserIdPrefix       = "cache:user:id:"
	cacheUserMobilePrefix   = "cache:user:mobile:"
	cacheUserUsernamePrefix = "cache:user:username:"
)

var group singleflight.Group

type UserModel interface {
	Create(ctx context.Context, data *User) error
	FindOne(ctx context.Context, id int64) (*User, error)
	FindByMobile(ctx context.Context, mobile string) (*[]User, error)
	FindByUsername(ctx context.Context, username string) (*[]User, error)
	Update(ctx context.Context, data *User) error
	Delete(ctx context.Context, id int64) error
}

type User struct {
	Id       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Mobile   string `db:"mobile"`
}

type UserExec struct {
	exec      *sql.Execute
	cacheExec *sql.CacheExecute
	logger    logger.Logger
}

func NewUserModel(db *gorm.DB, cache *sql.CacheExecute, logger logger.Logger) *UserExec {
	exec := sql.NewExecute(User{}, db)
	return &UserExec{
		exec:      exec,
		cacheExec: cache,
		logger:    logger,
	}
}

func (u *UserExec) Create(ctx context.Context, data *User) error {
	err := u.exec.Create(ctx, data)
	if err != nil {
		return err
	}
	go u.DelCache(ctx, data)
	return nil
}

func (u *UserExec) FindOne(ctx context.Context, id int64) (*User, error) {
	cachestr := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	result, err, _ := group.Do(cachestr, func() (interface{}, error) {
		datacache := u.Get(ctx, cachestr)
		if datacache != nil {
			return datacache, nil
		}
		var data User
		u.exec.AddWhere("id = ?", id)
		err := u.exec.Find(ctx, &data)
		if err != nil {
			return nil, err
		}
		go u.Set(cachestr, &data)
		return &data, nil
	})
	return result.(*User), err
}

func (u *UserExec) FindByMobile(ctx context.Context, mobile string) (*[]User, error) {
	cachestr := fmt.Sprintf("%s%v", cacheUserMobilePrefix, mobile)
	result, err, _ := group.Do(cachestr, func() (interface{}, error) {
		cacheval, err := u.cacheExec.GetCache(cachestr, ctx)
		datascache := u.GetMany(ctx, cacheval)
		if datascache != nil {
			return datascache, nil
		}
		var datas []User
		u.exec.AddWhere("mobile = ?", mobile)
		err = u.exec.Find(ctx, &datas)
		if err != nil {
			return nil, err
		}
		go u.SetMany(cachestr, &datas)
		return &datas, nil
	})
	return result.(*[]User), err
}

func (u *UserExec) FindByUsername(ctx context.Context, username string) (*[]User, error) {
	cachestr := fmt.Sprintf("%s%v", cacheUserUsernamePrefix, username)
	result, err, _ := group.Do(cachestr, func() (interface{}, error) {
		cacheval, err := u.cacheExec.GetCache(cachestr, ctx)
		datascache := u.GetMany(ctx, cacheval)
		if datascache != nil {
			return datascache, nil
		}
		var datas []User
		u.exec.AddWhere("username = ?", username)
		err = u.exec.Find(ctx, &datas)
		if err != nil {
			return nil, err
		}
		go u.SetMany(cachestr, &datas)
		return &datas, nil
	})
	return result.(*[]User), err
}

func (u *UserExec) Update(ctx context.Context, data *User) error {
	u.exec.AddWhere("id = ?", data.Id)
	err := u.exec.Update(ctx, data)
	if err != nil {
		return err
	}
	go u.DelCache(ctx, data)
	return nil
}

func (u *UserExec) Delete(ctx context.Context, id int64) error {
	var data User
	d, err := u.FindOne(ctx, id)
	if err != nil {
		return err
	}
	data = *d
	err = u.exec.Delete(ctx, &data)
	if err != nil {
		return err
	}
	go u.DelCache(ctx, &data)
	return nil
}

// 序列化
func UnMarshalJSON(s string, model *User) error {
	return json.Unmarshal([]byte(s), model)
}

func UnMarshalString(s string, model *[]int64) error {
	return json.Unmarshal([]byte(s), model)
}

// cache
func (u *UserExec) DelCache(ctx context.Context, model *User) {
	err := u.cacheExec.DeleteCache(fmt.Sprintf("%s%v", cacheUserIdPrefix, model.Id), ctx)
	if err != nil {
		u.logger.Error("Primary key cache delete failure, id = "+fmt.Sprintf("%v", model.Id), err)
	}
	err = u.cacheExec.DeleteCache(fmt.Sprintf("%s%v", cacheUserMobilePrefix, model.Mobile), ctx)
	if err != nil {
		u.logger.Warn("Non-primary key cache delete failure: ", err)
	}
	err = u.cacheExec.DeleteCache(fmt.Sprintf("%s%v", cacheUserUsernamePrefix, model.Username), ctx)
	if err != nil {
		u.logger.Warn("Non-primary key cache delete failure: ", err)
	}
}

func (u *UserExec) Get(ctx context.Context, cachestr string) *User {
	var data User
	cacheval, err := u.cacheExec.GetCache(cachestr, ctx)
	if err == nil {
		err := UnMarshalJSON(cacheval, &data)
		if err != nil {
			u.logger.Warn("UnMarshal failure: ", err)
			return nil
		}
		return &data
	}
	if !errors.Is(err, CacheNotFound) {
		u.logger.Warn("Primary key cache get failure: ", err)
		return nil
	}
	return nil
}

func (u *UserExec) GetMany(ctx context.Context, cachestr string) *[]User {
	var datas []User
	cacheval, err := u.cacheExec.GetCache(cachestr, ctx)
	if err == nil {
		var key []int64
		err := UnMarshalString(cacheval, &key)
		if err != nil {
			u.logger.Warn("UnMarshal failure: ", err)
			return nil
		}
		for _, c := range key {
			data, err := u.FindOne(ctx, c)
			if err != nil {
				return nil
			}
			datas = append(datas, *data)
		}
		return &datas
	}
	if !errors.Is(err, CacheNotFound) {
		u.logger.Warn("Primary key cache get failure: ", err)
		return nil
	}
	return nil
}

func (u *UserExec) Set(cachestr string, data *User) {
	ctx, cancel := context.WithTimeout(context.Background(), u.cacheExec.SetTTl)
	err := u.cacheExec.SetCache(cachestr, ctx, data)
	if err != nil {
		u.logger.Warn("Primary key cache set failure: ", err)
	}
	cancel()
}

func (u *UserExec) SetMany(cachestr string, data *[]User) {
	ctx, cancel := context.WithTimeout(context.Background(), u.cacheExec.SetTTl)
	var key []int64
	for _, v := range *data {
		key = append(key, v.Id)
		err := u.cacheExec.SetCache(fmt.Sprintf("%s%v", cacheUserIdPrefix, v.Id), ctx, &v)
		if err != nil {
			u.logger.Warn("Primary key cache set failure: ", err)
		}
	}
	err := u.cacheExec.SetCache(cachestr, ctx, &key)
	if err != nil {
		u.logger.Warn("Primary key cache set failure: ", err)
	}
	cancel()
}
