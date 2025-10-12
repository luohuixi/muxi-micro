package timeout

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http" // t_http
)

func setupRouterWithTimeout(handlerFunc gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(handlerFunc)

	r.GET("/normal", func(c *gin.Context) {
		c.JSON(http.StatusOK, t_http.FinalResp{
			Code:    0,
			Message: "OK",
			Data:    "success",
		})
	})

	r.GET("/sleep", func(c *gin.Context) {
		time.Sleep(2 * time.Second)
		c.JSON(http.StatusOK, t_http.FinalResp{
			Code:    0,
			Message: "OK",
			Data:    "should timeout",
		})
	})

	return r
}

func TestTimeoutMiddleware(t *testing.T) {

	t.Run("请求未超时", func(t *testing.T) {
		router := setupRouterWithTimeout(Timeout())
		req := httptest.NewRequest(http.MethodGet, "/normal", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"success"`)
	})

	t.Run("请求超时", func(t *testing.T) {
		router := setupRouterWithTimeout(Timeout(
			WithTimeoutDuration(1),
			WithTimeoutMessage("msg"),
			WithTimeoutCode(1),
		))
		req := httptest.NewRequest(http.MethodGet, "/sleep", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var resp t_http.FinalResp
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusGatewayTimeout, w.Code)
		assert.Equal(t, resp.Message, "msg")
		assert.Equal(t, resp.Code, 1)
	})
}
