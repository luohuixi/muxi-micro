package limiter

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/ulule/limiter/v3"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 返回一个只挂载了限流中间件的测试路由
func newRouter(opts ...Option) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Limiter(opts...))
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	return r
}

// 默认配置 + 非饱和流量：应全部 200
func TestLimiter_WithinQuota(t *testing.T) {
	r := newRouter(WithRate("5-S")) // 5 次/秒，测试方便
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("X-Forwarded-For", "127.0.0.1")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "pong", w.Body.String())
	}
}

// 多一次请求触发限流：应返回 429 & 自定义 JSON
func TestLimiter_RateLimited(t *testing.T) {
	r := newRouter(WithRate("2-S")) // 每秒 2 次

	// 连续三次请求：前两次 OK，第三次被限流
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("X-Forwarded-For", "192.0.2.1") // 固定同一 IP
		r.ServeHTTP(w, req)

		if i < 2 {
			assert.Equal(t, http.StatusOK, w.Code)
			continue
		}

		// 第 3 次：检查限流响应
		assert.Equal(t, http.StatusTooManyRequests, w.Code) // 429
		var resp t_http.FinalResp
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, DefaultCodeRateLimited, resp.Code)
		assert.Equal(t, DefaultMSGRateLimited, resp.Message)
	}
}

// 覆盖自定义选项
func TestLimiter_CustomOption(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		r := newRouter(
			WithRate("1-S"),
			WithCodeRateLimited(90001),
			WithMSGRateLimited("慢点儿~"),
		)

		// 第 1 次：通过
		w1 := httptest.NewRecorder()
		req1 := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req1.Header.Set("X-Forwarded-For", "203.0.113.1")
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// 第 2 次：立即再次请求 → 被限流
		w2 := httptest.NewRecorder()
		req2 := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req2.Header.Set("X-Forwarded-For", "203.0.113.1")
		r.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)

		var resp t_http.FinalResp
		_ = json.Unmarshal(w2.Body.Bytes(), &resp)
		assert.Equal(t, 90001, resp.Code)
		assert.Equal(t, "慢点儿~", resp.Message)

		// 等 1 秒后重试，应恢复可用
		time.Sleep(time.Second)
		w3 := httptest.NewRecorder()
		req3 := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req3.Header.Set("X-Forwarded-For", "203.0.113.1")
		r.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusOK, w3.Code)
	})

	t.Run("illegal rate string", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		newRouter(
			WithRate("~,>"),
			WithCodeRateLimited(90001),
			WithMSGRateLimited("慢点儿~"),
		)

	})

	t.Run("illegal rate string", func(t *testing.T) {

		r := newRouter(
			WithStore(&DummyStore{}),
			WithMSGRateError("出错了"),
		)

		// 应当失败
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("X-Forwarded-For", "203.0.113.1")
		r.ServeHTTP(w, req)

		var resp t_http.FinalResp
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 500, w.Code)
		assert.Equal(t, DefaultCodeRateError, resp.Code)
		assert.Equal(t, "出错了"+errNotImplemented.Error(), resp.Message)
	})

}

type DummyStore struct{}

var errNotImplemented = errors.New("not implemented")

func (s *DummyStore) Get(ctx context.Context, key string, rate limiter.Rate) (limiter.Context, error) {
	return limiter.Context{}, errNotImplemented
}

func (s *DummyStore) Peek(ctx context.Context, key string, rate limiter.Rate) (limiter.Context, error) {
	return limiter.Context{}, errNotImplemented
}

func (s *DummyStore) Reset(ctx context.Context, key string, rate limiter.Rate) (limiter.Context, error) {
	return limiter.Context{}, errNotImplemented
}

func (s *DummyStore) Increment(ctx context.Context, key string, count int64, rate limiter.Rate) (limiter.Context, error) {
	return limiter.Context{}, errNotImplemented
}
