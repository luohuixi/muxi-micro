package create

import (
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/muxi-Infra/muxi-micro/tool/gin/parse"
)

// 存在不覆盖
func CreateWire(addr, project string, apis []*parse.Api) error {
	dir := path.Join(addr, "wire.go")
	if _, err := os.Stat(dir); os.IsExist(err) {
		return nil
	}

	tmplPath := filepath.Join("gin", "template", "wire.tpl")

	t, err := template.New("wire").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := path.Join(addr, "wire.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	data := struct {
		Project string
		Name    []*Name
	}{
		Project: project,
		Name:    MaxFirstLetter(apis),
	}

	if err := t.ExecuteTemplate(file, "wire", data); err != nil {
		return err
	}

	return nil
}
