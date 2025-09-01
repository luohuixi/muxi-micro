package timeout

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
	"time"
)

const (
	DefaultTimeout = 10 * time.Second
	DefaultMessage = "请求超时，请稍后重试"
	CodeTimeout    = 10003
)

type timeoutCfg struct {
	duration time.Duration
	code     int
	message  string
}

type Option func(cfg *timeoutCfg)

func WithTimeoutDuration(d time.Duration) Option {
	return func(cfg *timeoutCfg) {
		cfg.duration = d
	}
}

func WithTimeoutCode(code int) Option {
	return func(cfg *timeoutCfg) {
		cfg.code = code
	}
}

func WithTimeoutMessage(msg string) Option {
	return func(cfg *timeoutCfg) {
		cfg.message = msg
	}
}

// 超时中间件
func Timeout(opts ...Option) gin.HandlerFunc {
	cfg := &timeoutCfg{
		duration: DefaultTimeout,
		code:     CodeTimeout,
		message:  DefaultMessage,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return timeout.New(
		timeout.WithTimeout(cfg.duration),

		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),

		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(504, t_http.CommonResp{
				Code:    cfg.code,
				Message: cfg.message,
				Data:    nil,
			})
		}),
	)
}
