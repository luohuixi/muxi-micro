package curd

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCurdError(t *testing.T) {
	t.Run("can not find model.go", func(t *testing.T) {
		tempDir, _ := ioutil.TempDir("", "curd_no_model_test")
		defer os.RemoveAll(tempDir)

		curdCmd := InitCurdCobra()

		var output strings.Builder
		curdCmd.SetOutput(&output)
		curdCmd.SetArgs([]string{
			"--dir", tempDir,
		})

		err := curdCmd.Execute()
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
	})

	t.Run("no primary key in model.go", func(t *testing.T) {
		tempDir, _ := ioutil.TempDir("", "curd_error_model_test")
		defer os.RemoveAll(tempDir)
		modelPath := filepath.Join(tempDir, "model.go")
		modelContent := `package model

type User struct {
	ID   int64   
	Name string ` + "`gorm:\"index\"`" + `
}
`
		_ = ioutil.WriteFile(modelPath, []byte(modelContent), 0644)
		var output strings.Builder
		curdCmd := InitCurdCobra()
		curdCmd.SetOutput(&output)
		curdCmd.SetArgs([]string{
			"--dir", tempDir,
		})

		err := curdCmd.Execute()
		assert.ErrorContains(t, err, "you should create a primary key")
	})

	t.Run("too many primary key in model.go", func(t *testing.T) {
		tempDir, _ := ioutil.TempDir("", "curd_error_model_test")
		defer os.RemoveAll(tempDir)
		modelPath := filepath.Join(tempDir, "model.go")
		modelContent := `package model

type User struct {
	ID   int64  ` + "`gorm:\"primaryKey;autoIncrement\"`" + `
	Name string ` + "`gorm:\"primaryKey;autoIncrement\"`" + `
}
`
		_ = ioutil.WriteFile(modelPath, []byte(modelContent), 0644)
		var output strings.Builder
		curdCmd := InitCurdCobra()
		curdCmd.SetOutput(&output)
		curdCmd.SetArgs([]string{
			"--dir", tempDir,
		})

		err := curdCmd.Execute()
		assert.ErrorContains(t, err, "only one primary key need to be created")
	})

	t.Run("not correct type primary key in model.go", func(t *testing.T) {
		tempDir, _ := ioutil.TempDir("", "curd_error_model_test")
		defer os.RemoveAll(tempDir)
		modelPath := filepath.Join(tempDir, "model.go")
		modelContent := `package model

type User struct {
	ID   int  ` + "`gorm:\"primaryKey;autoIncrement\"`" + `
	Name string
}
`
		_ = ioutil.WriteFile(modelPath, []byte(modelContent), 0644)
		var output strings.Builder
		curdCmd := InitCurdCobra()
		curdCmd.SetOutput(&output)
		curdCmd.SetArgs([]string{
			"--dir", tempDir,
		})

		err := curdCmd.Execute()
		assert.ErrorContains(t, err, "primary key type should be int64")
	})
}
