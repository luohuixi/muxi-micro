package create

import (
	"os"
	"path"
	"path/filepath"
	"text/template"
)

// type覆盖
func CreateType(addr, content, pkg string) error {
	tmplPath := filepath.Join("gin", "template", "type.tpl")

	t, err := template.New("type").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := path.Join(addr, "type.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	data := struct {
		PackName string
		Content  string
	}{
		PackName: pkg,
		Content:  content,
	}

	if err := t.ExecuteTemplate(file, "type", data); err != nil {
		return err
	}

	return nil
}
