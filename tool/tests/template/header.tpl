{{- define "header" -}}
package {{.PackageName}}

import (
	"testing"
)
{{template "test" $}}
{{- end -}}