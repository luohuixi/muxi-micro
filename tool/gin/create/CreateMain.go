package create

import (
	"os"
	"path"
	"path/filepath"
	"text/template"
)

// 存在不覆盖
func CreateMain(addr string) error {
	dir := path.Join(addr, "main.go")
	if _, err := os.Stat(dir); err == nil {
		return nil
	}

	tmplPath := filepath.Join("gin", "template", "main.tpl")

	t, err := template.New("main").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := path.Join(addr, "main.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	if err := t.ExecuteTemplate(file, "main", nil); err != nil {
		return err
	}

	return nil
}
