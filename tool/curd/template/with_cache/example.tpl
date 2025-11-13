{{- define "example" -}}
package {{.PackageName}}

import (
    "time"

	"github.com/muxi-Infra/muxi-micro/pkg/logger"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
)

var _ {{.ModelName}}Models = (*Extra{{.ModelName}}Exec)(nil)

type {{.ModelName}}Models interface {
	{{.ModelName}}Model
	// 可以在这里添加额外的方法接口
}

type Extra{{.ModelName}}Exec struct {
	*{{.ModelName}}Exec
	db    *gorm.DB
	cache *redis.Client
}

func New{{.ModelName}}Models(db *gorm.DB, cache *redis.Client, setTTL, expiration time.Duration, l logger.Logger) ({{.ModelName}}Models, error) {
	if err := db.AutoMigrate({{.ModelName}}{}); err != nil {
		return nil, err
	}

	instance := New{{.ModelName}}Model(db, cache, setTTL, expiration, l)

	return &Extra{{.ModelName}}Exec{
		instance,
		db,
		cache,
	}, nil
}

// 可以在这里添加额外的方法实现
// func (e *Extra{{.ModelName}}Exec) Example() {
//     e.db.Where().Find()
// }
{{- end -}}