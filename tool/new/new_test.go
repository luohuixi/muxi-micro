package new

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("dir not exist", func(t *testing.T) {
		newCmd := InitNewCobra()
		var output strings.Builder
		newCmd.SetOutput(&output)
		newCmd.SetArgs([]string{
			"--dir", "dir not exist",
		})
		err := newCmd.Execute()
		assert.ErrorContains(t, err, "can not find dir")
	})
}
