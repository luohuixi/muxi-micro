package create

import (
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/muxi-Infra/muxi-micro/tool/gin/parse"
)

// 存在不覆盖
func CreateLogic(addr, pkg string, service []*parse.Service) error {
	for _, s := range service {
		dir := path.Join(addr, s.Handler+".go")
		if _, err := os.Stat(dir); err == nil {
			continue
		}
		if s.Method.Req != "" {
			if err := CreateLogicWithReq(addr, pkg, s.Handler, s.Method.Req); err != nil {
				return err
			}
		} else {
			if err := CreateLogicNoReq(addr, pkg, s.Handler); err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateLogicWithReq(addr, pkg, handler, req string) error {
	tmplPath := filepath.Join("gin", "template", "logic_withreq.tpl")

	t, err := template.New("logic_withreq").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := path.Join(addr, handler+".go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	data := struct {
		PackName string
		Req      string
		Handler  string
	}{
		PackName: pkg,
		Req:      req,
		Handler:  handler,
	}

	if err := t.ExecuteTemplate(file, "logic_withreq", data); err != nil {
		return err
	}

	return nil
}

func CreateLogicNoReq(addr, pkg, handler string) error {
	tmplPath := filepath.Join("gin", "template", "logic_noreq.tpl")

	t, err := template.New("logic_noreq").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := path.Join(addr, handler+".go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	data := struct {
		PackName string
		Handler  string
	}{
		PackName: pkg,
		Handler:  handler,
	}

	if err := t.ExecuteTemplate(file, "logic_noreq", data); err != nil {
		return err
	}

	return nil
}
