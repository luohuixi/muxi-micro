{{- define "var" -}}
package {{.PackageName}}

import (
	"github.com/muxi-Infra/muxi-micro/pkg/sql"
)

var DBNotFound = sql.DBNotFound
{{- end -}}