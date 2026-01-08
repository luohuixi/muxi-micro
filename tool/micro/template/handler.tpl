{{- define "handler" -}}

package handler

import (
    "YourPath/{{.Path}}"
)

type {{.ServiceName}} struct {
	{{.Pkg}}.Unimplemented{{.ServiceName}}Server
}

{{- end -}}