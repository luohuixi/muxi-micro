{{- define "var" -}}
package {{.PackageName}}

import (
	"github.com/muxi-Infra/muxi-micro/pkg/sql"
)

const CacheNotFound = sql.CacheNotFound

var DBNotFound = sql.DBNotFound
{{- end -}}