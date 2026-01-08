{{- define "main" -}}

package main

import (
    "log"
    "context"
    "YourPath/handler"
    "YourPath/server"
)

func main() {
     // TODO: 替换成注入后的结构体
     s, err := server.{{.ServiceName}}Server(new(handler.{{.ServiceName}}))
     if err != nil {
         log.Fatal(err)
     }

     ctx := context.Background()
     if err := s.Serve(ctx); err != nil {
         log.Fatal(err)
     }
}

{{- end -}}