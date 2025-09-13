package new

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("dir not exist", func(t *testing.T) {
		newCmd := InitNewCobra()
		var output strings.Builder
		newCmd.SetOutput(&output)
		newCmd.SetArgs([]string{
			"--dir", "test",
		})
		err := newCmd.Execute()
		assert.ErrorContains(t, err, "The system cannot find the path specified")
	})
}
