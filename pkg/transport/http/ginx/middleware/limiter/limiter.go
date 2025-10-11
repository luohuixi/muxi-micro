package limiter

import (
	"github.com/gin-gonic/gin"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/handler"
	"github.com/ulule/limiter/v3"
	l_gin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"net/http"
)

const (
	DefaultRateStr         = "200-S"
	DefaultCodeRateLimited = 42901
	DefaultCodeRateError   = 50001
	DefaultMSGRateLimited  = "请求太频繁，请稍后再试"
	DefaultMSGRateError    = "限流器出错:"
)

type limiterCfg struct {
	rateStr         string
	store           limiter.Store
	codeRateLimited int
	msgRateLimited  string
	codeRateError   int
	msgRateError    string
}

type Option func(cfg *limiterCfg)

func WithRate(rateStr string) Option {
	return func(cfg *limiterCfg) {
		cfg.rateStr = rateStr
	}
}

func WithStore(store limiter.Store) Option {
	return func(cfg *limiterCfg) {
		cfg.store = store
	}
}

func WithCodeRateLimited(code int) Option {
	return func(cfg *limiterCfg) {
		cfg.codeRateLimited = code
	}
}

func WithMSGRateError(message string) Option {
	return func(cfg *limiterCfg) {
		cfg.msgRateError = message
	}
}

func WithMSGRateLimited(message string) Option {
	return func(cfg *limiterCfg) {
		cfg.msgRateLimited = message
	}
}

// 限流中间件
func Limiter(opts ...Option) gin.HandlerFunc {
	var cfg = &limiterCfg{
		rateStr:         DefaultRateStr,
		codeRateLimited: DefaultCodeRateLimited,
		codeRateError:   DefaultCodeRateError,
		store:           memory.NewStore(),
		msgRateLimited:  DefaultMSGRateLimited,
		msgRateError:    DefaultMSGRateError,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	rate, err := limiter.NewRateFromFormatted(cfg.rateStr)
	if err != nil {
		panic(err)
	}

	lim := limiter.New(cfg.store, rate)
	// 自定义限流返回结构
	return l_gin.NewMiddleware(lim,
		l_gin.WithLimitReachedHandler(func(c *gin.Context) {
			handler.HandleResponse(c, t_http.Response{
				HttpCode: http.StatusTooManyRequests,
				Code:     cfg.codeRateLimited,
				Message:  cfg.msgRateLimited,
				Data:     nil,
			})
		}),

		l_gin.WithErrorHandler(func(c *gin.Context, err error) {
			handler.HandleResponse(c, t_http.Response{
				HttpCode: http.StatusInternalServerError,
				Code:     cfg.codeRateError,
				Message:  cfg.msgRateError + err.Error(),
				Data:     nil,
			})
		}),
	)
}
