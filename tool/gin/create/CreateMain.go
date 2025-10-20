package create

import (
	"os"
	"path"
	"path/filepath"
	"text/template"
)

// 存在不覆盖
func CreateMain(addr, project string) error {
	dir := path.Join(addr, "main.go")
	if _, err := os.Stat(dir); os.IsExist(err) {
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

	data := struct {
		Project string
	}{
		Project: project,
	}

	if err := t.ExecuteTemplate(file, "main", data); err != nil {
		return err
	}

	return nil
}
