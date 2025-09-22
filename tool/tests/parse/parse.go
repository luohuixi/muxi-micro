package parse

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

var TooManyErr = errors.New("to many error are not recommended")

type Value struct {
	Param     string
	ParamType string
}

type FuncStruct struct {
	Receive  string
	FuncName string
	Params   []Value
	Returns  []Value
}

func ParsePackage(filePath string) (string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}
	return node.Name.Name, nil
}

func ParseFunc(filePath string) (*[]FuncStruct, *[]FuncStruct, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, nil, err
	}

	var Funcs, Reces []FuncStruct
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			var Func FuncStruct
			var flag bool
			if fn.Recv != nil {
				Func.Receive = exprToString(fn.Recv.List[0].Type)
			}
			Func.FuncName = fn.Name.Name
			if fn.Type.Params != nil {
				for _, field := range fn.Type.Params.List {
					typeStr := exprToString(field.Type)
					for _, name := range field.Names {
						Func.Params = append(Func.Params, Value{
							Param:     identToString(name),
							ParamType: typeStr,
						})
					}
				}
			}
			if fn.Type.Results != nil {
				for _, field := range fn.Type.Results.List {
					typeStr := exprToString(field.Type)
					var str string
					if typeStr == "error" {
						if flag == true {
							return nil, nil, TooManyErr
						}
						str = "err"
						flag = true
					} else {
						str = "expected" + fmt.Sprintf("%v", len(Func.Returns)+1)
					}
					Func.Returns = append(Func.Returns, Value{
						Param:     str,
						ParamType: typeStr,
					})
				}
			}
			if Func.Receive == "" {
				Funcs = append(Funcs, Func)
			} else {
				Reces = append(Reces, Func)
			}
		}
	}
	return &Funcs, &Reces, nil
}

func identToString(ident *ast.Ident) string {
	if ident == nil {
		return ""
	}
	return ident.Name
}

func exprToString(expr ast.Expr) string {
	switch v := expr.(type) {
	// 基础类型
	case *ast.Ident:
		return v.Name
	// 包（如: context.Context）
	case *ast.SelectorExpr:
		return exprToString(v.X) + "." + v.Sel.Name
	// 指针
	case *ast.StarExpr:
		return "*" + exprToString(v.X)
	// 切片或数组
	case *ast.ArrayType:
		if v.Len == nil {
			return "[]" + exprToString(v.Elt)
		}
		return "[" + exprToString(v.Len) + "]" + exprToString(v.Elt)
	// 哈希表
	case *ast.MapType:
		return "map[" + exprToString(v.Key) + "]" + exprToString(v.Value)
	// 可变参数（如: ...int）
	case *ast.Ellipsis:
		return "..." + exprToString(v.Elt)
	// 管道
	case *ast.ChanType:
		switch v.Dir {
		case ast.SEND:
			return "chan<- " + exprToString(v.Value)
		case ast.RECV:
			return "<-chan " + exprToString(v.Value)
		default:
			return "chan " + exprToString(v.Value)
		}
	// 否则 interface 或 struct，func 等未命名的返回值返回 interface
	default:
		return "interface{}"
	}
}
