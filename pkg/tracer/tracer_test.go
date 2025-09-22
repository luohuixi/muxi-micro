package tracer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomError(t *testing.T) {
	t.Run("not correct random1", func(t *testing.T) {
		_, err := NewZipkin(
			"http://localhost:9411/api/v2/spans",
			"demo_service",
			"localhost:50051",
			-1.0,
		)
		assert.EqualError(t, err, "random number must be between 0 and 1")
	})
	t.Run("not correct random2", func(t *testing.T) {
		_, err := NewJaeger(
			"http://localhost:14268/api/traces",
			"demo_service",
			2.0,
		)
		assert.EqualError(t, err, "random number must be between 0 and 1")
	})
	t.Run("not correct random3", func(t *testing.T) {
		_, err := NewSkyWalking(
			"localhost:11800",
			"demo_service",
			"demo_instance",
			100,
		)
		assert.EqualError(t, err, "random number must be between 0 and 1")
	})
}
