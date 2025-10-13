package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTest(t *testing.T) {
	t.Run("file not write", func(t *testing.T) {
		testCmd := InitTestCobra()
		var output strings.Builder
		testCmd.SetOut(&output)
		err := testCmd.Execute()
		assert.ErrorContains(t, err, "\"file\" not set")
	})

	t.Run("file not exist", func(t *testing.T) {
		testCmd := InitTestCobra()
		var output strings.Builder
		testCmd.SetOut(&output)
		testCmd.SetArgs([]string{
			"--file", "nonexist.go",
		})
		err := testCmd.Execute()
		assert.Error(t, err)
		assert.True(t,
			strings.Contains(err.Error(), "The system cannot find the file specified") ||
				strings.Contains(err.Error(), "no such file or directory"),
			"错误消息不匹配",
		)
	})
}
