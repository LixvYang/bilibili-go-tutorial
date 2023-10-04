package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lixvyang/ginzero"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

func GinXid(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		xid := xid.New().String()
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("xid", xid)
		})

		c.Header("xid", xid)

		c.Set("logger", logger)
		c.Next()
	}
}

func NewLogger() zerolog.Logger {
	logger := zerolog.
		New(os.Stdout).
		With().
		Caller().
		Timestamp().
		Logger()
	return logger
}

func main() {
	logger := NewLogger()
	r := gin.New()

	r.Use(GinXid(logger), ginzero.Ginzero(&logger), ginzero.RecoveryWithZero(&logger, true))

	r.GET("/hello", func(c *gin.Context) {
		zlog := c.MustGet("logger").(zerolog.Logger)

		zlog.Info().Msg("发生了 Inf o级别的日志消息")
		zlog.Debug().Msg("发生了 Debug 级别的日志消息")
		zlog.Error().Msg("发生了 Error 级别的日志消息")

		c.String(200, "hello")
	})

	r.GET("/panic", func(c *gin.Context) {
		panic("panic msg.")
	})

	r.Run(":8002")
}
