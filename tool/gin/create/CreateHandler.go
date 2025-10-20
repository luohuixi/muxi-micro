package create

import (
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/muxi-Infra/muxi-micro/tool/gin/parse"
)

// 会覆盖
func CreateHandler(addr, pkg string, service []*parse.Service, server *parse.Server) error {
	tmplPath := filepath.Join("gin", "template", "handler.tpl")

	t, err := template.New("handler").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	pkg2 := maxFirstLetter(pkg)
	outputPath := path.Join(addr, "handler.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	data := struct {
		PackName string
		Name     string
		Service  []*parse.Service
		Prefix   string
		Group    string
	}{
		PackName: pkg,
		Name:     pkg2,
		Service:  service,
		Prefix:   server.Prefix,
		Group:    server.Group,
	}

	if err := t.ExecuteTemplate(file, "handler", data); err != nil {
		return err
	}

	return nil
}
