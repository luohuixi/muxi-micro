{{- define "main" -}}
package main

import (
	"{{.Project}}/router"
)

func main() {
	if err := router.Run("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}

{{- end -}}