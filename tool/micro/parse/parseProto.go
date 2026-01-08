package parse

import (
	"errors"
	"path"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"
)

type ProtoRpc struct {
	Method   string
	Request  string
	Response string
}

type ProtoFile struct {
	ServiceName string
	Rpc         []*ProtoRpc
	Path        string
	Pkg         string
}

func ParseProto(protoPath string) (*ProtoFile, error) {
	parser := &protoparse.Parser{
		ImportPaths: []string{filepath.Dir(protoPath)},
	}

	files, err := parser.ParseFiles(filepath.Base(protoPath))
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.New("no proto files found")
	}

	fd := files[0]
	// 规定一个 proto 一个 service
	if len(fd.GetServices()) > 1 {
		return nil, errors.New("only one service in one proto supported")
	}

	var protoFile ProtoFile
	for _, s := range fd.GetServices() {
		protoFile.ServiceName = s.GetName()
		for _, method := range s.GetMethods() {
			var rpc ProtoRpc
			rpc.Method = method.GetName()
			rpc.Request = method.GetInputType().GetName()
			rpc.Response = method.GetOutputType().GetName()

			protoFile.Rpc = append(protoFile.Rpc, &rpc)
		}
	}

	goImportPath, goPkgName := parseGoPackage(fd.GetFileOptions().GetGoPackage())
	protoFile.Path = path.Clean(goImportPath)
	protoFile.Pkg = goPkgName
	return &protoFile, nil
}

func parseGoPackage(goPackage string) (string, string) {
	var goImportPath, goPkgName string

	if idx := strings.LastIndex(goPackage, ";"); idx != -1 {
		goImportPath = goPackage[:idx]
		goPkgName = goPackage[idx+1:]
	} else {
		goImportPath = goPackage
	}

	if goPkgName == "" {
		goPkgName = path.Base(goImportPath)
	}

	return goImportPath, goPkgName
}
