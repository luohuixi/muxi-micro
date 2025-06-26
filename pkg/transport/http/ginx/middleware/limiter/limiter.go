package limiter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
	"github.com/ulule/limiter/v3"
	l_gin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

const (
	DefaultRateStr         = "1000-s"
	DefaultCodeRateLimited = 42901
	DefaultCodeRateError   = 50001
)

type limiterCfg struct {
	rateStr         string
	store           limiter.Store
	codeRateLimited int
	codeRateError   int
}

type LimiterOption func(cfg *limiterCfg)

func WithRate(rateStr string) LimiterOption {
	return func(cfg *limiterCfg) {
		cfg.rateStr = rateStr
	}
}

func WithStore(store limiter.Store) LimiterOption {
	return func(cfg *limiterCfg) {
		cfg.store = store
	}
}

func WithCodeRateLimited(code int) LimiterOption {
	return func(cfg *limiterCfg) {
		cfg.codeRateLimited = code
	}
}

func WithCodeRateError(code int) LimiterOption {
	return func(cfg *limiterCfg) {
		cfg.codeRateError = code
	}
}

// 限流中间件
func Limiter(opts ...LimiterOption) gin.HandlerFunc {
	var cfg = &limiterCfg{
		rateStr:         DefaultRateStr,
		codeRateLimited: DefaultCodeRateLimited,
		codeRateError:   DefaultCodeRateError,
		store:           memory.NewStore(),
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
			c.AbortWithStatusJSON(429, t_http.CommonResp{
				Message: "请求太频繁，请稍后再试",
				Code:    cfg.codeRateLimited,
				Data:    nil,
			},
			)
		}),

		l_gin.WithErrorHandler(func(c *gin.Context, err error) {
			c.AbortWithStatusJSON(500, t_http.CommonResp{
				Message: fmt.Sprintf("限流器出错:", err.Error()),
				Code:    cfg.codeRateError,
				Data:    nil,
			})
		}),
	)
}
