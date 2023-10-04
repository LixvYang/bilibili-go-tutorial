package main

import (
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func GinLogger(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		cost := time.Since(start)
		defer logger.Info().
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Str("ip", c.ClientIP()).
			Str("user-agent", c.Request.UserAgent()).
			Str("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()).
			Dur("cost", cost).Send()

		c.Next()
	}
}

// GinRecovery recovers possible project panic
func GinRecovery(logger zerolog.Logger, stack bool) gin.HandlerFunc {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.
						Error().
						Str("path", c.Request.URL.Path).
						Any("error", err).
						Str("request", string(httpRequest)).
						Send()

					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					errors.New(string(debug.Stack()))
					logger.
						Error().
						Stack().
						Err(errors.New(string(debug.Stack()))).
						Str("error", "[Recovery from panic]").
						Str("request", string(httpRequest)).
						Send()

				} else {
					logger.
						Error().
						Str("error", "[Recovery from panic]").
						Any("error", err).
						Str("request", string(httpRequest)).
						Send()
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

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
