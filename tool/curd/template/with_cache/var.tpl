{{- define "var" -}}
package {{.PackageName}}

import (
    "github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/sql"
	"time"
)

type CacheStruct struct {
	RedisAddr     string
	RedisPassword string
	Number        int
	TtlForCache   time.Duration
	TtlForSet     time.Duration
	Log           logger.Logger
}

const CacheNotFound = sql.CacheNotFound

var DBNotFound = sql.DBNotFound
{{- end -}}