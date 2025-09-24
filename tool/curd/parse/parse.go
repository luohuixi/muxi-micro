package parse

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
)

var (
	ErrNoPrimaryKey   = errors.New("you should create a primary key")
	ErrManyPrimaryKey = errors.New("only one primary key need to be created")
	ErrPrimaryKeyType = errors.New("primary key type should be int64")
)

type FieldInfo struct {
	Name string // 结构体字段名
	Type string // 字段类型
	Tag  string // 原始标签
}

type StructInfo struct {
	Name    string // 结构体名
	Primary []FieldInfo
	Index   []FieldInfo
}

// ParseStructs 解析文件中的所有结构体并返回它们的信息
func ParseStructs(filePath string) ([]StructInfo, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	structs, err := parseFile(content)
	if err != nil {
		return nil, err
	}

	// 对每个结构体进行检查
	for _, s := range structs {
		if len(s.Primary) == 0 {
			return nil, fmt.Errorf("%w in struct %s", ErrNoPrimaryKey, s.Name)
		}
		if len(s.Primary) > 1 {
			return nil, fmt.Errorf("%w in struct %s", ErrManyPrimaryKey, s.Name)
		}
		if s.Primary[0].Type != "int64" {
			return nil, fmt.Errorf("%w in struct %s", ErrPrimaryKeyType, s.Name)
		}
	}

	return structs, nil
}

// parseFile 解析Go文件中的所有结构体
func parseFile(content []byte) ([]StructInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var structs []StructInfo

	ast.Inspect(f, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		info := StructInfo{
			Name: typeSpec.Name.Name,
		}

		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 || field.Tag == nil {
				continue
			}

			fieldName := field.Names[0].Name
			fieldType := getTypeName(field.Type)

			fieldInfo := FieldInfo{
				Name: fieldName,
				Type: fieldType,
			}

			tag := strings.Trim(field.Tag.Value, "`")
			fieldInfo.Tag = getTagValue(tag, "gorm")
			if fieldInfo.Tag != "" {
				if strings.Contains(fieldInfo.Tag, "primaryKey") {
					info.Primary = append(info.Primary, fieldInfo)
				}
				if strings.Contains(fieldInfo.Tag, "index") || strings.Contains(fieldInfo.Tag, "unique") {
					info.Index = append(info.Index, fieldInfo)
				}
			}
		}

		structs = append(structs, info)
		return true
	})

	return structs, nil
}

func getTypeName(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	default:
		return "unknown"
	}
}

func getTagValue(tag, key string) string {
	tags := reflect.StructTag(tag)
	return tags.Get(key)
}
