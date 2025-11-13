{{- define "transaction" -}}
package {{.PackageName}}

import (
	"context"

	"gorm.io/gorm"
)

type TransactionExec struct {
	db *gorm.DB
	// TODO: 添加需要事务处理的interface
	// 例子
	// u UserModels
	// o OrderModels
}

func NewTranExec(db *gorm.DB) (*TransactionExec, error) {
    // TODO: 添加对应事务的实例
	// u, _ := NewUserModels(db)
	return &TransactionExec{
		db: db,
		//u: u,
		//o: o,
	}, nil
}

func Transaction(ctx context.Context, db *gorm.DB, fn func(context.Context, *TransactionExec) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		t, err := NewTranExec(tx)
		if err != nil {
		    return err
		}

		return fn(ctx, t)
	})
}
{{- end -}}