{{- define "header" -}}
package {{.PackageName}}

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
    {{- range $field := .Fields}}
    cache{{$.ModelName}}{{$field}}Prefix = "cache:{{$.ModelName}}:{{$field}}:"
    {{- end}}
)

var group singleflight.Group

type {{.ModelName}}Model interface {
    Create(ctx context.Context, data *{{.ModelName}}) error
    FindOne(ctx context.Context, id int64) (*{{.ModelName}}, error)
    {{- range $notPr := .NotPrs}}
    FindBy{{$notPr}}(ctx context.Context, {{$notPr}} string) (*[]{{$.ModelName}}, error)
    {{- end}}
    Update(ctx context.Context, data *{{.ModelName}}) error
    Delete(ctx context.Context, id int64) error
}

type {{.ModelName}}Exec struct {
    exec      *sql.Execute
    cacheExec *sql.CacheExecute
    logger    logger.Logger
}

func New{{.ModelName}}Model(db *gorm.DB, cache *sql.CacheExecute, logger logger.Logger) *{{.ModelName}}Exec {
    exec := sql.NewExecute({{.ModelName}}{}, db)
    return &{{.ModelName}}Exec{
        exec:      exec,
        cacheExec: cache,
        logger:    logger,
    }
}

// 序列化
func UnMarshalJSON(s string, model *{{.ModelName}}) error {
    return json.Unmarshal([]byte(s), model)
}

func UnMarshalString(s string, model *[]int64) error {
    return json.Unmarshal([]byte(s), model)
}

{{template "db" $}}

{{template "cache" $}}
{{- end -}}