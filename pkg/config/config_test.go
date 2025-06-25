package config

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

		c, _ := LoadFromLocal(tmpFile.Name())
		m := make(map[string]string, 10)
		key := "database.mysql.host"

		value, _ := GetConfig(c, key)
		m[key] = value
		assert.Equal(t, "127.0.0.1", m[key])
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

		_, err := LoadFromLocal(tmpFile.Name())
		assert.EqualError(t, err, "only .yaml, .yml, or .json are supported")
	})

	t.Run("change data", func(t *testing.T) {
		tmpFile, _ := os.CreateTemp("", "test-watch-*.yaml")
		defer os.Remove(tmpFile.Name())

		initialContent := `
          database:
            host: "localhost"
            port: 3306
        `
		_, _ = tmpFile.WriteString(initialContent)
		tmpFile.Close()

		v, _ := LoadFromLocal(tmpFile.Name())
		m := make(map[string]string, 10)
		key := "database.mysql.host"

		value, _ := GetConfig(v, key)
		m[key] = value

		ctx, cancel := context.WithCancel(context.Background())

		configMap := map[string]string{"database.host": "localhost"}
		updateCh := WatchConfig(v, configMap, ctx)
		time.Sleep(2 * time.Second) // 确保监听就绪了

		newContent := `
           database:
             host: "192.168.1.100"
             port: 3306
        `
		_ = os.WriteFile(tmpFile.Name(), []byte(newContent), 0644)

		go func() {
			time.Sleep(5 * time.Second)
			cancel()
		}()

		<-updateCh
		assert.Equal(t, "192.168.1.100", configMap["database.host"])
	})

	t.Run("file not exist", func(t *testing.T) {
		_, err := LoadFromLocal("nonexistent.yaml")
		assert.EqualError(t, err, "open nonexistent.yaml: The system cannot find the file specified.")
	})
}