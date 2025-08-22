package create

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

func safeFilename(tableName string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, tableName)
}

func CreateVar(pkg, dir string) error {
	tmplPath := filepath.Join("curd", "template", "var.tpl")

	t, err := template.New("var").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(dir, "var.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	data := struct {
		PackageName string
	}{
		PackageName: pkg,
	}

	if err := t.ExecuteTemplate(file, "var", data); err != nil {
		return err
	}

	return nil
}

func CreateExample(pkg, dir, table string) error {
	tmplPath := filepath.Join("curd", "template", "example.tpl")

	t, err := template.New("example").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(dir, safeFilename(table)+"Model.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	data := struct {
		PackageName string
		ModelName   string
	}{
		PackageName: pkg,
		ModelName:   table,
	}

	if err := t.ExecuteTemplate(file, "example", data); err != nil {
		return err
	}

	return nil
}

func CreateExample_gen(pkg, dir, table string, fields []string) error {
	tmplPath := []string{
		filepath.Join("curd", "template", "header.tpl"),
		filepath.Join("curd", "template", "cache.tpl"),
		filepath.Join("curd", "template", "db.tpl"),
	}

	t, err := template.New("header").ParseFiles(tmplPath...)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(dir, safeFilename(table)+"Model_gen.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	data := struct {
		PackageName string
		ModelName   string
		Fields      []string
		NotPrs      []string
		Pr          string
	}{
		PackageName: pkg,
		ModelName:   table,
		Fields:      fields,
		NotPrs:      fields[:len(fields)-1],
		Pr:          fields[len(fields)-1],
	}

	if err := t.ExecuteTemplate(file, "header", data); err != nil {
		return err
	}

	// 设置为只读权限
	if err := file.Chmod(0444); err != nil {
		return err
	}
	return nil
}
