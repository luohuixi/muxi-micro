{{- define "transaction" -}}
package {{.PackageName}}

import (
	"context"

	"github.com/muxi-Infra/muxi-micro/pkg/sql"
	"gorm.io/gorm"
)

type TranExec struct {
	exec *sql.Execute
	// TODO: 添加需要事务处理的interface
	// 例子
	// u UserModels
	// o OrderModels
}

func NewTranExec(dsn string) (*TranExec, error) {
	db, err := sql.ConnectDB(dsn, nil)
	if err != nil {
		return nil, err
	}
	exec := sql.NewExecute(db)
	return &TranExec{
    	exec,
    	//nil,
    	//nil,
    }, nil
}

func (t *TranExec) Transaction(ctx context.Context, fn func(context.Context, *TranExec) error) error {
	return t.exec.Transaction(func(tx *gorm.DB) error {
		exec := sql.NewExecute(tx)
		//userexec := &UserExec{
        //	exec: exec,
        //}
        //orderexec := &OrderExec{
        //	exec: exec,
        //}
		tran := &TranExec{
			exec: exec,
			//u: &ExtraUserExec{
            //	UserExec: userexec,
            //	db: tx,
            //},
            //o: &ExtraOrderExec{
            //	OrderExec: orderexec,
            //	db: tx,
            //},
		}
		return fn(ctx, tran)
	})
}
{{- end -}}