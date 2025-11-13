{{- define "example" -}}
package {{.PackageName}}

import (
	"gorm.io/gorm"
)

var _ {{.ModelName}}Models = (*Extra{{.ModelName}}Exec)(nil)

type {{.ModelName}}Models interface {
	{{.ModelName}}Model
	// 可以在这里添加额外的方法接口
}

type Extra{{.ModelName}}Exec struct {
	*{{.ModelName}}Exec
	db *gorm.DB
}

func New{{.ModelName}}Models(db *gorm.DB) ({{.ModelName}}Models, error) {
	if err := db.AutoMigrate({{.ModelName}}{}); err != nil {
    	return nil, err
    }
    instance := New{{.ModelName}}Model(db)

    return &Extra{{.ModelName}}Exec{
    	instance,
    	db,
    }, nil
}

// 可以在这里添加额外的方法实现
// func (e *Extra{{.ModelName}}Exec) Example() {
//     e.db.Where().Find()
// }
{{- end -}}