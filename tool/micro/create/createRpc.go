package create

import (
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/muxi-Infra/muxi-micro/tool/micro/parse"
)

func CreateRpc(cover bool, output string, protoFile *parse.ProtoFile) error {
	tmplPath := filepath.Join("micro", "template", "rpc.tpl")

	t, err := template.New("rpc").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := path.Join(output, "rpc")

	if err := os.MkdirAll(path.Dir(outputPath), 0755); err != nil {
		return err
	}

	for _, r := range protoFile.Rpc {
		outputPath := path.Join(output, "handler", r.Method+".go")
		if _, err := os.Stat(outputPath); err == nil && !cover {
			continue
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		data := struct {
			ServiceName string
			Rpc         *parse.ProtoRpc
			Path        string
			Pkg         string
		}{
			ServiceName: protoFile.ServiceName,
			Rpc:         r,
			Path:        protoFile.Path,
			Pkg:         protoFile.Pkg,
		}

		if err := t.ExecuteTemplate(file, "rpc", data); err != nil {
			return err
		}
	}
	return nil
}
