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

func New{{.ModelName}}Models(DBdsn string,  Cache *CacheStruct) ({{.ModelName}}Models, error) {
	db, err := sql.ConnectDB(DBdsn, {{.ModelName}}{})
	if err != nil {
		return nil, err
	}
	cache := sql.ConnectCache(Cache.RedisAddr, Cache.RedisPassword, Cache.Number, Cache.TtlForCache, Cache.TtlForSet)

	instance := New{{.ModelName}}Model(db, cache, Cache.Log)

	return &Extra{{.ModelName}}Exec{
		instance,
	}, nil
}

// 可以在这里添加额外的方法实现
// func (e *Extra{{.ModelName}}Exec) ...
{{- end -}}