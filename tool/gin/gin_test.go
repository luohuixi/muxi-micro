package gin

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestGinError(t *testing.T) {
	t.Run("can not find xxx.api", func(t *testing.T) {
		tempDir, _ := os.MkdirTemp("", "curd_no_model_test")
		filePath := path.Join(tempDir, "xxx.api")
		_, err := os.Open(filePath)

		ginCmd := InitGinCobra()

		var output strings.Builder
		ginCmd.SetOut(&output)
		ginCmd.SetArgs([]string{
			"--api", tempDir,
		})

		err = ginCmd.Execute()
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
	})
}
