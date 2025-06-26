package config

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

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
		yamlContent := `
          database:
            mysql:
               host: "127.0.0.1"
               port: "3306"
               username: "root"
               password: "password"
        `
		_, _ = tmpFile.WriteString(yamlContent)
		tmpFile.Close()

		cfg, _ := NewLocalConfig(&Config{}, "config.yaml")
		cfg.LoadData()
		assert.Equal(t, "127.0.0.1", cfg.GetData().(*Config).Database.MySQL.Host)
	})

	t.Run("not correct type", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp("", "db-config-*.txt")
		defer os.Remove(tmpFile.Name())
		yamlContent := `
           database:
              mysql:
              host: "127.0.0.1"
              port: "3306"
              username: "root"
              password: "password"
        `
		_, _ = tmpFile.WriteString(yamlContent)
		tmpFile.Close()

		cfg, _ := NewLocalConfig(&Config{}, tmpFile.Name())
		cfg.LoadData()
		assert.EqualError(t, cfg.GetErr(), "only .yaml, .yml, or .json are supported")
	})

	t.Run("file not exist", func(t *testing.T) {
		cfg, _ := NewLocalConfig(&Config{}, "nonexistent.yaml")
		cfg.LoadData()
		assert.EqualError(t, cfg.GetErr(), "open nonexistent.yaml: The system cannot find the file specified.")
	})
}
