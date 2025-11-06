package sql

import (
	"context"
	"errors"
	"reflect"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrNonPointer = errors.New("data must be a pointer")
	DBNotFound    = errors.New("this data is empty")
)

// model为自动迁移，可不设置
func ConnectDB(dsn string, model any) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}

	if model != nil {
		if err := db.AutoMigrate(model); err != nil {
			return nil, err
		}
	}

	return db, nil
}

type Execute struct {
	Model *gorm.DB
}

func NewExecute(db *gorm.DB) *Execute {
	return &Execute{
		Model: db,
	}
}

// 执行部分
func (e *Execute) Create(ctx context.Context, data, model any) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return ErrNonPointer
	}
	return e.Model.Model(model).WithContext(ctx).Create(data).Error
}

func (e *Execute) Update(ctx context.Context, data, model any) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return ErrNonPointer
	}
	return e.Model.Model(model).WithContext(ctx).Save(data).Error
}

func (e *Execute) Delete(ctx context.Context, data, model any) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return ErrNonPointer
	}
	return e.Model.Model(model).WithContext(ctx).Delete(data).Error
}

func (e *Execute) Find(ctx context.Context, data, model any) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return ErrNonPointer
	}
	return e.Model.Model(model).WithContext(ctx).Find(data).Error
}

func (e *Execute) Transaction(fn func(*gorm.DB) error) error {
	return e.Model.Transaction(fn)
}
