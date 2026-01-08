package create

import (
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/muxi-Infra/muxi-micro/tool/micro/parse"
)

func CreateServer(cover bool, output string, protoFile *parse.ProtoFile) error {
	tmplPath := filepath.Join("micro", "template", "server.tpl")

	t, err := template.New("server").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := path.Join(output, "server/server.go")

	if _, err := os.Stat(outputPath); err == nil && !cover {
		return nil
	}

	if err := os.MkdirAll(path.Dir(outputPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := t.ExecuteTemplate(file, "server", protoFile); err != nil {
		return err
	}
	return nil
}
