package cors

// TODO 测试可用性，没弄明白怎么用代码测试跨域问题,现在的
import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCors_Default(t *testing.T) {
	t.Run("if is nil", func(t *testing.T) {
		assert.NotEmpty(t, Cors())
	})
}

func TestCors_CustomOption(t *testing.T) {
	t.Run("if is nil", func(t *testing.T) {
		c := Cors(
			WithCorsOrigins([]string{"*"}),
			WithCorsAllowCredentials(true),
			WithCorsAllowMethods("GET", "POST"),
			WithCorsMaxAge(10*time.Second),
			WithCorsAllowHeaders("Content-Type", "Accept"),
		)
		assert.NotEmpty(t, c)
	})
}
