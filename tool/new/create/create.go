package create

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var DirNotEsxitErr = errors.New("can not find dir")

func CreateDocument(rootDir string) error {
	dirs := []string{
		filepath.Join(rootDir, "helloworld", "api", "local", "v1"),
		filepath.Join(rootDir, "helloworld", "api", "third_party"),
		filepath.Join(rootDir, "helloworld", "cmd"),
		filepath.Join(rootDir, "helloworld", "configs"),
		filepath.Join(rootDir, "helloworld", "internal", "infrastructure"),
		filepath.Join(rootDir, "helloworld", "internal", "repository"),
		filepath.Join(rootDir, "helloworld", "internal", "server"),
		filepath.Join(rootDir, "helloworld", "internal", "service"),
		filepath.Join(rootDir, "helloworld", "internal", "wire"),
	}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateCmd(rootDir string) error {
	path := []string{"main.go"}
	for _, p := range path {
		templatePath := filepath.Join("new", "template", "cmd", p+".tpl")
		templateContent, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}
		createPath := filepath.Join(rootDir, "helloworld", "cmd", p)
		if err := ioutil.WriteFile(createPath, templateContent, 0644); err != nil {
			return err
		}
	}
	return nil
}

func CreateConfigs(rootDir string) error {
	path := []string{"config.yaml", "load.go"}
	for _, p := range path {
		templatePath := filepath.Join("new", "template", "configs", p+".tpl")
		templateContent, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}
		createPath := filepath.Join(rootDir, "helloworld", "configs", p)
		if err := ioutil.WriteFile(createPath, templateContent, 0644); err != nil {
			return err
		}
	}
	return nil
}

func CreateInfrastructure(rootDir string) error {
	path := []string{"interceptor.go", "provider.go"}
	for _, p := range path {
		templatePath := filepath.Join("new", "template", "infrastructure", p+".tpl")
		templateContent, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}
		createPath := filepath.Join(rootDir, "helloworld", "internal", "infrastructure", p)
		if err := ioutil.WriteFile(createPath, templateContent, 0644); err != nil {
			fmt.Println(666)
			return err
		}
	}
	return nil
}

func CreateLocal(rootDir string) error {
	path := []string{"hello.pb.go", "hello.proto", "hello_grpc.pb.go"}
	for _, p := range path {
		templatePath := filepath.Join("new", "template", "local", p+".tpl")
		templateContent, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}
		createPath := filepath.Join(rootDir, "helloworld", "api", "local", "v1", p)
		if err := ioutil.WriteFile(createPath, templateContent, 0644); err != nil {
			return err
		}
	}
	return nil
}

func CreateRepository(rootDir string) error {
	paths := []string{"model.go", "provider.go", "Usermodel.go", "Usermodel_gen.go", "var.go"}
	for _, p := range paths {
		templatePath := filepath.Join("new", "template", "repository", p+".tpl")
		templateContent, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}
		createPath := filepath.Join(rootDir, "helloworld", "internal", "repository", p)
		if err := ioutil.WriteFile(createPath, templateContent, 0644); err != nil {
			return err
		}
	}
	return nil
}

func CreateServer(rootDir string) error {
	paths := []string{"provider.go", "register.go"}
	for _, p := range paths {
		templatePath := filepath.Join("new", "template", "server", p+".tpl")
		templateContent, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}
		createPath := filepath.Join(rootDir, "helloworld", "internal", "server", p)
		if err := ioutil.WriteFile(createPath, templateContent, 0644); err != nil {
			return err
		}
	}
	return nil
}

func CreateService(rootDir string) error {
	paths := []string{"hello.go", "provider.go"}
	for _, p := range paths {
		templatePath := filepath.Join("new", "template", "service", p+".tpl")
		templateContent, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}
		createPath := filepath.Join(rootDir, "helloworld", "internal", "service", p)
		if err := ioutil.WriteFile(createPath, templateContent, 0644); err != nil {
			return err
		}
	}
	return nil
}

func CreateWire(rootDir string) error {
	paths := []string{"wire.go", "wire_gen.go"}
	for _, p := range paths {
		templatePath := filepath.Join("new", "template", "wire", p+".tpl")
		templateContent, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}
		createPath := filepath.Join(rootDir, "helloworld", "internal", "wire", p)
		if err := ioutil.WriteFile(createPath, templateContent, 0644); err != nil {
			return err
		}
	}
	return nil
}

func CreateExplain(rootDir string) error {
	paths := []string{"explain.md", "go.mod", "go.sum", "Dockerfile"}
	for _, p := range paths {
		templatePath := filepath.Join("new", "template", p+".tpl")
		templateContent, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}
		createPath := filepath.Join(rootDir, "helloworld", p)
		if err := ioutil.WriteFile(createPath, templateContent, 0644); err != nil {
			return err
		}
	}
	return nil
}

func CreateAll(rootDir string) error {
	if _, err := os.Stat(rootDir); err != nil {
		return DirNotEsxitErr
	}
	if err := CreateDocument(rootDir); err != nil {
		return err
	}
	if err := CreateCmd(rootDir); err != nil {
		return err
	}
	if err := CreateConfigs(rootDir); err != nil {
		return err
	}
	if err := CreateInfrastructure(rootDir); err != nil {
		return err
	}
	if err := CreateLocal(rootDir); err != nil {
		return err
	}
	if err := CreateRepository(rootDir); err != nil {
		return err
	}
	if err := CreateServer(rootDir); err != nil {
		return err
	}
	if err := CreateService(rootDir); err != nil {
		return err
	}
	if err := CreateWire(rootDir); err != nil {
		return err
	}
	if err := CreateExplain(rootDir); err != nil {
		return err
	}
	return nil
}
