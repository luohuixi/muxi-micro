package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var YAMLTEST = `
database:
  mysql:
    host: "127.0.0.1"
    port: "3306"
    username: "root"
    password: "password"

  redis:
    host: "127.0.0.1"
    port: "6379"
    username: "root"
    passowrd: "password"
    db: 0
`

type DatabaseConfig struct {
	MySQL struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mysql"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

func TestLoadFromLocal(t *testing.T) {
	t.Run("correct read", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp("", "db-config-*.yaml")
		defer os.Remove(tmpFile.Name())
		_, _ = tmpFile.WriteString(YAMLTEST)
		tmpFile.Close()

		cfg, _ := NewLocalConfig(&Config{}, tmpFile.Name())
		_ = cfg.LoadData()
		assert.Equal(t, "127.0.0.1", cfg.GetData().(*Config).Database.MySQL.Host)
	})

	t.Run("not correct type", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp("", "db-config-*.txt")
		defer os.Remove(tmpFile.Name())
		_, _ = tmpFile.WriteString(YAMLTEST)
		tmpFile.Close()

		cfg, _ := NewLocalConfig(&Config{}, tmpFile.Name())
		err := cfg.LoadData()
		assert.EqualError(t, err, "only .yaml, .yml, or .json are supported")
	})

	t.Run("file not exist", func(t *testing.T) {
		cfg, _ := NewLocalConfig(&Config{}, "nonexistent.yaml")
		err := cfg.LoadData()
		assert.EqualError(t, err, "open nonexistent.yaml: The system cannot find the file specified.")
	})

	//nacos不好测
}
