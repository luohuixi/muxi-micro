package log

import (
	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/logx"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

// TestGenLogID 测试生成的 LogID 格式和唯一性
func TestGenLogID(t *testing.T) {
	id1 := genLogID("test")
	id2 := genLogID("test")

	if id1 == id2 {
		t.Errorf("expected unique logID, got same: %s", id1)
	}

	// 匹配前缀 test-xxxxxx
	re := regexp.MustCompile(`^test-[0-9a-f]{16}$`)
	if !re.MatchString(id1) {
		t.Errorf("invalid logID format: %s", id1)
	}
}

// TestSetAndGetLogID 测试在 Gin Context 中正确设置与获取
func TestSetAndGetLogID(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/", nil)
	logID := "abc-123"
	SetLogID(ctx, logID)

	got := GetLogID(ctx)
	if got != logID {
		t.Errorf("expected %s, got %s", logID, got)
	}
}

// TestAutoGenerateLogID 当没有 LogID 时自动生成
func TestAutoGenerateLogID(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("GET", "/", nil)
	SetLogID(ctx, genLogID(GetGlobalName(ctx)))
	id := GetLogID(ctx)
	if id == "" {
		t.Errorf("expected auto-generated logID, got empty")
	}
	if !regexp.MustCompile(`^unknown-[0-9a-f]{16}$`).MatchString(id) {
		t.Errorf("unexpected auto logID format: %s", id)
	}
}

// TestGlobalMiddleware 测试中间件能正确注入 logID 与 logger
func TestGlobalMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	stdLogger := logx.NewStdLogger()

	router.Use(GlobalNameMiddleware("mockService"))
	router.Use(GlobalLoggerMiddleware(stdLogger))

	router.GET("/ping", func(ctx *gin.Context) {
		l := GetLogger(ctx)
		if l == nil {
			t.Errorf("logger not injected into context")
		}

		logID := GetLogID(ctx)
		if logID == "" {
			t.Errorf("logID should not be empty")
		}

		// 记录一条日志，确保不会 panic
		l.Info("middleware test", logger.Field{"req": "ping"})

		ctx.JSON(http.StatusOK, gin.H{
			"ok":    true,
			"logID": logID,
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	t.Logf("response: %s", w.Body.String())
}

// TestGetGlobalName 测试全局服务名逻辑
func TestGetGlobalName(t *testing.T) {
	ctx := &gin.Context{}

	// 未设置时默认返回 unknown
	if name := GetGlobalName(ctx); name != DefaultName {
		t.Errorf("expected %s, got %s", DefaultName, name)
	}

	// 设置后应返回对应值
	GlobalNameMiddleware("testService")(ctx)
	name := GetGlobalName(ctx)
	if name != "testService" {
		t.Errorf("expected testService, got %s", name)
	}
}

// TestLogIDPerformance 简单测试生成 ID 的性能
func TestLogIDPerformance(t *testing.T) {
	start := time.Now()
	for i := 0; i < 10000; i++ {
		_ = genLogID("perf")
	}
	elapsed := time.Since(start)
	t.Logf("generated 10k logIDs in %v", elapsed)
}
