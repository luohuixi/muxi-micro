package create

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/muxi-Infra/muxi-micro/tool/gin/parse"
)

// 会覆盖
func Create2Service(addr, pkg string, service []*parse.Service) error {
	tmplPath := filepath.Join("gin", "template", "service.tpl")

	t, err := template.New("service").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	pkg2 := maxFirstLetter(pkg)
	outputPath := path.Join(addr, pkg+".go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	data := struct {
		PackName string
		Name     string
		Service  []*parse.Service
	}{
		PackName: pkg,
		Name:     pkg2,
		Service:  service,
	}

	if err := t.ExecuteTemplate(file, "service", data); err != nil {
		return err
	}

	return nil
}

func maxFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}
