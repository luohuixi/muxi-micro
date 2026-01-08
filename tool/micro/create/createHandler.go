package create

import (
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/muxi-Infra/muxi-micro/tool/micro/parse"
)

func CreateHandler(cover bool, output string, protoFile *parse.ProtoFile) error {
	tmplPath := filepath.Join("micro", "template", "handler.tpl")

	t, err := template.New("handler").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := path.Join(output, "handler/handler.go")

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

	if err := t.ExecuteTemplate(file, "handler", protoFile); err != nil {
		return err
	}
	return nil
}
