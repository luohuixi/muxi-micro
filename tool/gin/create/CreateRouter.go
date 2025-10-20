package create

import (
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/muxi-Infra/muxi-micro/tool/gin/parse"
)

type Name struct {
	Max string
	Min string
}

// 存在不覆盖
func CreateRouter(addr, project string, apis []*parse.Api) error {
	dir := path.Join(addr, "router.go")
	if _, err := os.Stat(dir); os.IsExist(err) {
		return nil
	}

	tmplPath := filepath.Join("gin", "template", "router.tpl")

	t, err := template.New("router").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := path.Join(addr, "router.go")
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

	if err := t.ExecuteTemplate(file, "router", data); err != nil {
		return err
	}

	return nil
}

func MaxFirstLetter(apis []*parse.Api) []*Name {
	var names []*Name
	for _, api := range apis {
		var name Name
		name.Max = maxFirstLetter(api.ServiceName)
		name.Min = api.ServiceName
		names = append(names, &name)
	}
	return names
}
