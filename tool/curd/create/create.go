package create

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/muxi-Infra/muxi-micro/tool/curd/parse"
	"gorm.io/gorm/schema"
)

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

func CreateExample_gen(pkg, dir, table string, fields []parse.FieldInfo, open bool) error {
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

	gormNames := gormName(fields)

	data := struct {
		PackageName string
		ModelName   string
		Fields      []parse.FieldInfo
		NotPrs      []parse.FieldInfo
		Pr          string
		GNotPrs     []string
		GPr         string
	}{
		PackageName: pkg,
		ModelName:   table,
		Fields:      fields,
		NotPrs:      fields[:len(fields)-1],
		Pr:          fields[len(fields)-1].Name,
		GNotPrs:     gormNames[:len(gormNames)-1],
		GPr:         gormNames[len(gormNames)-1],
	}

	if err := t.ExecuteTemplate(file, "header", data); err != nil {
		return err
	}

	return nil
}

func CreateCache(pkg, dir, table string, fields []parse.FieldInfo, open bool) error {
	if !open {
		return nil
	}
	tmplPath := filepath.Join("curd", "template", "with_cache", "cache.tpl")

	t, err := template.New("cache").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(dir, safeFilename(table)+"cache.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()

	data := struct {
		PackageName string
		ModelName   string
		NotPrs      []parse.FieldInfo
		Pr          string
	}{
		PackageName: pkg,
		ModelName:   table,
		NotPrs:      fields[:len(fields)-1],
		Pr:          fields[len(fields)-1].Name,
	}

	if err := t.ExecuteTemplate(file, "cache", data); err != nil {
		return err
	}

	return nil
}

func CreateTranscation(pkg, dir string, transcation bool) error {
	if transcation == false {
		return nil
	}
	var tmplPath string
	tmplPath = filepath.Join("curd", "template", "transaction.tpl")
	t, err := template.New("header").ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(dir, "transaction.go")
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

	if err := t.ExecuteTemplate(file, "transaction", data); err != nil {
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

// 自动迁移特殊字段处理
func gormName(f []parse.FieldInfo) []string {
	var gormNames []string
	namingStrategy := schema.NamingStrategy{}

	for _, fi := range f {
		gormNames = append(gormNames, namingStrategy.ColumnName("", fi.Name))
	}

	return gormNames
}
