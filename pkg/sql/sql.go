package sql

import (
	"context"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"reflect"
)

var ErrNonPointer = errors.New("data must be a pointer")

func ConnectDB(dsn string, model any) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(model); err != nil {
		return nil, err
	}

	return db, nil
}

type Execute struct {
	QueryOptions []func(db *gorm.DB) *gorm.DB
	Model        *gorm.DB
}

func NewExecute(model any, db *gorm.DB) *Execute {
	return &Execute{
		QueryOptions: nil,
		Model:        db.Model(model),
	}
}

// 拼接查询部分
func (e *Execute) AddWhere(str string, val any) {
	e.QueryOptions = append(e.QueryOptions, func(db *gorm.DB) *gorm.DB {
		return db.Where(str, val)
	})
}

func (e *Execute) Build(ctx context.Context) *gorm.DB {
	query := e.Model.WithContext(ctx)
	for _, opt := range e.QueryOptions {
		query = opt(query)
	}
	e.QueryOptions = nil
	return query
}

// 可补充查询类型

// 执行部分
func (e *Execute) Create(ctx context.Context, data any) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return ErrNonPointer
	}
	return e.Build(ctx).Create(data).Error
}

func (e *Execute) Update(ctx context.Context, data any) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return ErrNonPointer
	}
	return e.Build(ctx).Save(data).Error
}

func (e *Execute) Delete(ctx context.Context, data any) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return ErrNonPointer
	}
	return e.Build(ctx).Delete(data).Error
}

func (e *Execute) Find(ctx context.Context, data any) error {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		return ErrNonPointer
	}
	return e.Build(ctx).Find(data).Error
}
