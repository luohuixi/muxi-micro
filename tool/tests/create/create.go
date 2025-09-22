package create

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/muxi-Infra/muxi-micro/tool/tests/parse"
)

type ParseStruct struct {
	PackageName string
	Func        *[]parse.FuncStruct
}

func CreateFunc(filePath string, funcStructs *[]parse.FuncStruct, packName string) error {
	if len(*funcStructs) == 0 {
		return nil
	}
	var tmplPath []string
	tmplPath = []string{
		filepath.Join("tests", "template", "header.tpl"),
		filepath.Join("tests", "template", "test.tpl"),
	}

	t, err := template.New("example").ParseFiles(tmplPath...)
	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)
	outputPath := filepath.Join(dir, GetFileName(filePath)+"_test.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := &ParseStruct{
		PackageName: packName,
		Func:        funcStructs,
	}
	if err := t.ExecuteTemplate(file, "header", data); err != nil {
		return err
	}

	return nil
}

func CreateRece(filePath string, receStructs *[]parse.FuncStruct, packName string) error {
	if len(*receStructs) == 0 {
		return nil
	}
	var tmplPath []string
	tmplPath = []string{
		filepath.Join("tests", "template", "header.tpl"),
		filepath.Join("tests", "template", "injectTest.tpl"),
	}

	t, err := template.New("example").ParseFiles(tmplPath...)
	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)
	outputPath := filepath.Join(dir, GetFileName(filePath)+"_test2.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := &ParseStruct{
		PackageName: packName,
		Func:        receStructs,
	}
	if err := t.ExecuteTemplate(file, "header", data); err != nil {
		return err
	}

	return nil
}

func GetFileName(filePath string) string {
	filename := filepath.Base(filePath)
	return strings.TrimSuffix(filename, ".go")
}
