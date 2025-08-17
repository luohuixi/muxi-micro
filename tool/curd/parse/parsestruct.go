package parse

import (
	"errors"
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

func ParseStruct(filename string) (string, []string, error) {
	filePath := "model/model.go"
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, err
	}

	structInfos, err := parseFile(content)
	if err != nil {
		return "", nil, err
	}

	if len(structInfos.Primary) == 0 {
		return "", nil, ErrNoPrimaryKey
	}
	if len(structInfos.Primary) > 1 {
		return "", nil, ErrManyPrimaryKey
	}
	if structInfos.Primary[0].Type != "int64" {
		return "", nil, ErrPrimaryKeyType
	}

	var index []string
	for _, structInfo := range structInfos.Index {
		if structInfo.Type == "unknown" {
			continue
		}
		index = append(index, structInfo.Name)
	}

	return structInfos.Primary[0].Name, index, nil
}

// parseFile 解析Go文件中的结构体
func parseFile(content []byte) (*StructInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var structInfos StructInfo

	ast.Inspect(f, func(n ast.Node) bool {
		// 只处理type
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		// 只处理struct
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}
		info := StructInfo{
			Name: typeSpec.Name.Name,
		}

		for _, field := range structType.Fields.List {
			// 跳过嵌入字段
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
				if isPrimaryKey := strings.Contains(fieldInfo.Tag, "primaryKey"); isPrimaryKey {
					info.Primary = append(info.Primary, fieldInfo)
				}
				if isIndex := strings.Contains(fieldInfo.Tag, "index") || strings.Contains(fieldInfo.Tag, "unique"); isIndex {
					info.Index = append(info.Index, fieldInfo)
				}
			}
		}

		structInfos = info
		return true
	})

	return &structInfos, nil
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
