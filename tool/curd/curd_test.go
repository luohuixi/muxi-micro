package curd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurdError(t *testing.T) {
	t.Run("can not find model.go", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "curd_no_model_test")
		defer os.RemoveAll(tempDir)

		curdCmd := InitCurdCobra()

		var output strings.Builder
		curdCmd.SetOut(&output)
		curdCmd.SetArgs([]string{
			"--dir", tempDir,
		})

		err := curdCmd.Execute()
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
	})

	t.Run("no primary key in model.go", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "curd_error_model_test")
		defer os.RemoveAll(tempDir)
		modelPath := filepath.Join(tempDir, "model.go")
		modelContent := `package model

type User struct {
	ID   int64   
	Name string ` + "`gorm:\"index\"`" + `
}
`
		_ = os.WriteFile(modelPath, []byte(modelContent), 0644)
		var output strings.Builder
		curdCmd := InitCurdCobra()
		curdCmd.SetOut(&output)
		curdCmd.SetArgs([]string{
			"--dir", tempDir,
		})

		err := curdCmd.Execute()
		assert.ErrorContains(t, err, "you should create a primary key")
	})

	t.Run("too many primary key in model.go", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "curd_error_model_test")
		defer os.RemoveAll(tempDir)
		modelPath := filepath.Join(tempDir, "model.go")
		modelContent := `package model

type User struct {
	ID   int64  ` + "`gorm:\"primaryKey;autoIncrement\"`" + `
	Name string ` + "`gorm:\"primaryKey;autoIncrement\"`" + `
}
`
		_ = os.WriteFile(modelPath, []byte(modelContent), 0644)
		var output strings.Builder
		curdCmd := InitCurdCobra()
		curdCmd.SetOut(&output)
		curdCmd.SetArgs([]string{
			"--dir", tempDir,
		})

		err := curdCmd.Execute()
		assert.ErrorContains(t, err, "only one primary key need to be created")
	})

	t.Run("not correct type primary key in model.go", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "curd_error_model_test")
		defer os.RemoveAll(tempDir)
		modelPath := filepath.Join(tempDir, "model.go")
		modelContent := `package model

type User struct {
	ID   int  ` + "`gorm:\"primaryKey;autoIncrement\"`" + `
	Name string
}
`
		_ = os.WriteFile(modelPath, []byte(modelContent), 0644)
		var output strings.Builder
		curdCmd := InitCurdCobra()
		curdCmd.SetOut(&output)
		curdCmd.SetArgs([]string{
			"--dir", tempDir,
		})

		err := curdCmd.Execute()
		assert.ErrorContains(t, err, "primary key type should be int64")
	})
}
