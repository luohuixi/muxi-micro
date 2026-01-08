{{- define "rpc" -}}

package handler

import (
    "YourPath/{{.Path}}"
    "context"
)

func (s *{{.ServiceName}}) {{.Rpc.Method}}(ctx context.Context, req *{{.Pkg}}.{{.Rpc.Request}}) (resp *{{.Pkg}}.{{.Rpc.Response}},err error) {
	// TODO: 实现具体逻辑
	return
}

{{- end -}}
