package ginx

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
	"github.com/stretchr/testify/assert"
)

type demoReq struct {
	Name string `json:"name" form:"name" binding:"required"`
}

type demoClaims struct {
	UID int64 `json:"uid"`
}

func handleDemoReq(c *gin.Context, r demoReq) t_http.Response {
	return t_http.Response{
		HttpCode: http.StatusOK,
		CommonResp: t_http.CommonResp{
			Code:    0,
			Message: "ok",
			Data:    r,
		},
	}
}

func handleWithClaims(c *gin.Context, uc demoClaims) t_http.Response {
	return t_http.Response{
		HttpCode: http.StatusOK,
		CommonResp: t_http.CommonResp{
			Code:    0,
			Message: "ok",
			Data:    uc,
		},
	}
}

func handleWithClaimsAndReq(c *gin.Context, r demoReq, uc demoClaims) t_http.Response {
	return t_http.Response{
		HttpCode: http.StatusOK,
		CommonResp: t_http.CommonResp{
			Code:    0,
			Message: "ok",
			Data: gin.H{
				"name": r.Name,
				"uid":  uc.UID,
			},
		},
	}
}

func handleNoBody(c *gin.Context) t_http.Response {
	return t_http.Response{
		HttpCode:   http.StatusOK,
		CommonResp: t_http.CommonResp{Code: 0, Message: "pong"},
	}
}

func mockClaimsOK(c *gin.Context) (demoClaims, error) {
	return demoClaims{UID: 1}, nil
}

var authError = errors.New("认证失败")

func mockClaimsFail(c *gin.Context) (demoClaims, error) {
	return demoClaims{}, authError
}

func decodeResp[T any](body []byte) (T, error) {
	var resp t_http.CommonResp
	var t T
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return t, err
	}
	db, _ := json.Marshal(resp.Data)
	err = json.Unmarshal(db, &t)
	return t, err
}

func setErrorHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Error(errors.New("test"))
	}
}

func TestWrapReq(t *testing.T) {
	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.POST("/demo", WrapReq(handleDemoReq))

	tests := []struct {
		name   string
		body   string
		status int
		expect string
	}{
		{"非法参数", `{"foo":"bar"}`, http.StatusBadRequest, ""},
		{"正常参数", `{"name":"bar"}`, http.StatusOK, "bar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/demo", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			g.ServeHTTP(w, req)

			assert.Equal(t, tt.status, w.Code)
			if tt.status == http.StatusOK {
				res, err := decodeResp[demoReq](w.Body.Bytes())
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, res.Name)
			}
		})
	}

	t.Run("error in front Middleware", func(t *testing.T) {
		g.GET("/testFrontError", setErrorHandlerFunc(), WrapReq(handleDemoReq))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/testFrontError", bytes.NewBufferString("{\"name\":\"bar\"}"))
		req.Header.Set("Content-Type", "application/json")
		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

	})
}

func TestWrapClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.GET("/ok", WrapClaims(mockClaimsOK, handleWithClaims))
	g.GET("/fail", WrapClaims(mockClaimsFail, handleWithClaims))

	t.Run("验证失败", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/fail", nil)
		g.ServeHTTP(w, req)

		var resp t_http.CommonResp
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "登陆状态异常:"+authError.Error(), resp.Message)
	})

	t.Run("验证成功", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/ok", nil)
		g.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error in front Middleware", func(t *testing.T) {
		g.GET("/testFrontError", setErrorHandlerFunc(), WrapClaims(mockClaimsOK, handleWithClaims))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/testFrontError", bytes.NewBufferString("{\"name\":\"bar\"}"))
		req.Header.Set("Content-Type", "application/json")
		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

	})
}

func TestWrapClaimsAndReq(t *testing.T) {
	gin.SetMode(gin.TestMode)
	g := gin.New()
	g.POST("/full", WrapClaimsAndReq(mockClaimsOK, handleWithClaimsAndReq))
	g.GET("/full", WrapClaimsAndReq(mockClaimsOK, handleWithClaimsAndReq))
	g.POST("/fail", WrapClaimsAndReq(mockClaimsFail, handleWithClaimsAndReq))

	t.Run("Success get", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/full?name=bob", nil)
		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp t_http.CommonResp
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		db, _ := json.Marshal(resp.Data)
		var m map[string]any
		_ = json.Unmarshal(db, &m)

		assert.Equal(t, "ok", resp.Message)
		assert.Equal(t, "bob", m["name"])
		assert.Equal(t, float64(1), m["uid"]) // json.Unmarshal => float64
	})

	t.Run("Success post", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/full", bytes.NewBufferString(`{"name":"bob"}`))
		req.Header.Set("Content-Type", "application/json")

		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp t_http.CommonResp
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		db, _ := json.Marshal(resp.Data)
		var m map[string]any
		_ = json.Unmarshal(db, &m)

		assert.Equal(t, "ok", resp.Message)
		assert.Equal(t, "bob", m["name"])
		assert.Equal(t, float64(1), m["uid"])
	})

	t.Run("BindErr", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/full", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json")
		g.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("AuthErr", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/fail", bytes.NewBufferString(`{"name":"alice"}`))
		req.Header.Set("Content-Type", "application/json")
		g.ServeHTTP(w, req)

		var resp t_http.CommonResp
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "登陆状态异常:"+authError.Error(), resp.Message)
	})

	t.Run("error in front Middleware", func(t *testing.T) {
		g.GET("/testFrontError", setErrorHandlerFunc(), WrapClaimsAndReq(mockClaimsOK, handleWithClaimsAndReq))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/testFrontError", bytes.NewBufferString("{\"name\":\"bar\"}"))
		req.Header.Set("Content-Type", "application/json")
		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

	})
}

func TestWrap(t *testing.T) {
	gin.SetMode(gin.TestMode)
	g := gin.New()

	t.Run("请求正常", func(t *testing.T) {
		g.GET("/ping", Wrap(handleNoBody))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp t_http.CommonResp
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "pong", resp.Message)
	})

	t.Run("error in front Middleware", func(t *testing.T) {
		g.GET("/testFrontError", setErrorHandlerFunc(), Wrap(handleNoBody))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/testFrontError", bytes.NewBufferString("{\"name\":\"bar\"}"))
		req.Header.Set("Content-Type", "application/json")
		g.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

	})

}

func TestDefaultEngine(t *testing.T) {
	t.Run("pprof enabled", func(t *testing.T) {
		g := NewDefaultEngine(
			WithEngine(gin.Default()),
			WithEnv(t_http.EnvTest),
		)
		g.GET("/ping", Wrap(handleNoBody))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/debug/pprof/", nil)
		g.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("pprof disabled", func(t *testing.T) {
		g := NewDefaultEngine(
			WithEngine(gin.Default()),
			WithEnv(t_http.EnvProd),
		)
		SetBindErrCode(42062)
		SetGetClaimsErrCode(12345)
		UseDefaultMiddleware(g)
		g.GET("/ping", Wrap(handleNoBody))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/debug/pprof/", nil)
		g.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
