package create

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

func CreateVar(pkg, dir string, open, cover bool) error {
	if cover == false && !CheckExist(dir, "var.go") {
		return nil
	}
	var tmplPath string
	if open {
		tmplPath = filepath.Join("curd", "template", "with_cache", "var.tpl")
	} else {
		tmplPath = filepath.Join("curd", "template", "no_cache", "var.tpl")
	}

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

func CreateExample(pkg, dir, table string, open, cover bool) error {
	if cover == false && !CheckExist(dir, safeFilename(table)+"model.go") {
		return nil
	}
	var tmplPath string
	if open {
		tmplPath = filepath.Join("curd", "template", "with_cache", "example.tpl")
	} else {
		tmplPath = filepath.Join("curd", "template", "no_cache", "example.tpl")
	}

	t, err := template.New("example").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(dir, safeFilename(table)+"model.go")
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

func CreateExample_gen(pkg, dir, table string, fields []string, open bool) error {
	var tmplPath []string
	if open {
		tmplPath = []string{
			filepath.Join("curd", "template", "with_cache", "header.tpl"),
			filepath.Join("curd", "template", "with_cache", "cache.tpl"),
			filepath.Join("curd", "template", "with_cache", "db.tpl"),
		}
	} else {
		tmplPath = []string{
			filepath.Join("curd", "template", "no_cache", "header.tpl"),
			filepath.Join("curd", "template", "no_cache", "db.tpl"),
		}
	}

	t, err := template.New("header").ParseFiles(tmplPath...)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(dir, safeFilename(table)+"model_gen.go")
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

	return nil
}

func safeFilename(tableName string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, tableName)
}

func CheckExist(dir, filename string) bool {
	path := filepath.Join(dir, filename)
	if _, err := os.Stat(path); err != nil {
		return true
	}
	return false
}
