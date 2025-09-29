package ginx

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"time"
)

const (
	LogIDKey    = "Gin-LogID"
	LoggerKey   = "Gin-Logger"
	NameKey     = "Gin-Name"
	DefaultName = "unknown"
)

// 生成日志 ID
func genLogID(prefix string) string {
	// 当前时间纳秒 + 随机字节混合
	timeBytes := []byte(fmt.Sprintf("%d", time.Now().UnixNano()))

	//生成随机短id
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	randomBytes := hex.EncodeToString(b)
	combined := append(timeBytes, randomBytes...)

	// SHA1 哈希处理
	hash := sha1.Sum(combined)
	shortHash := hex.EncodeToString(hash[:8]) // 取前8字节（16个字符）

	logID := fmt.Sprintf("%s-%s", prefix, shortHash)
	return logID
}

// 设置到响应中需要
func SetLogID(ctx *gin.Context, logID string) {
	ctx.Set(LogIDKey, logID)
}

// 设置到响应中需要
func GetLogID(ctx *gin.Context) string {
	var logID string
	value, ok := ctx.Get(LogIDKey)
	if !ok {
		// 先尝试从X-Request-ID请求头里面取
		logID = ctx.Request.Header.Get("X-Request-ID")
		// 如果不存在则尝试去生成一个
		if logID == "" {
			logID = genLogID(GetGlobalName(ctx))
		}
		SetLogID(ctx, logID)
		return logID
	}
	return value.(string)
}

// 用于设置在上下文中获取携带了特殊信息的日志,主动打印需要
func SetLogger(ctx *gin.Context, logger logger.Logger) {
	ctx.Set(LoggerKey, logger)
}

// 用于获取在上下文中获取携带了特殊信息的日志
func GetLogger(ctx *gin.Context) logger.Logger {
	ginLogger, ok := ctx.Get(LoggerKey)
	if !ok {
		return nil // 如果不存在需要返回一个空指针
	}
	return ginLogger.(logger.Logger)
}

func GlobalLoggerMiddleware(l logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		l = l.With(logger.Field{
			"logID": GetLogID(ctx), // 保证ctx中的所有日志都是自带logID的
		})
		SetLogger(ctx, l)
	}
}

func GlobalNameMiddleware(name string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(NameKey, name)
	}
}

func GetGlobalName(ctx *gin.Context) string {
	value, ok := ctx.Get(NameKey)
	if !ok {
		return DefaultName
	}
	return value.(string)
}
