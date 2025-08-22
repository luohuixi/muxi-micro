{{- define "example" -}}
package {{.PackageName}}

import (
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/sql"
	"time"
)

var _ {{.ModelName}}Models = (*Extra{{.ModelName}}Exec)(nil)

type {{.ModelName}}Models interface {
	{{.ModelName}}Model
	// 可以在这里添加额外的方法接口
}

type Extra{{.ModelName}}Exec struct {
	*{{.ModelName}}Exec
}

func New{{.ModelName}}Models(DBdsn, redisAddr, redisPassword string, number int, ttlForCache, ttlForSet time.Duration, l logger.Logger) ({{.ModelName}}Models, error) {
	db, err := sql.ConnectDB(DBdsn, {{.ModelName}}{})
	if err != nil {
		return nil, err
	}
	cache := sql.ConnectCache(redisAddr, redisPassword, number, ttlForCache, ttlForSet)

	instance := New{{.ModelName}}Model(db, cache, l)

	return &Extra{{.ModelName}}Exec{
		instance,
	}, nil
}

// 可以在这里添加额外的方法实现
// func (e *Extra{{.ModelName}}Exec) ...
{{- end -}}