package server

import (
	"time"

	"blockpropeller.dev/lib/log"
	"github.com/labstack/echo"
)

// LoggerMiddleware returns a middleware that logs HTTP requests.
func LoggerMiddleware(logger log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			logFields := log.Fields{
				"method":     req.Method,
				"uri":        req.RequestURI,
				"path":       req.URL.Path,
				"status":     res.Status,
				"latency_ms": stop.Sub(start).Milliseconds(),
			}
			if err != nil {
				logFields["error"] = err.Error()
			}
			intErr := c.Get("_internal_error")
			if intErr != nil {
				logFields["internal_error"] = intErr.(error).Error()
			}

			logger.Info("HTTP Request", logFields)

			return nil
		}
	}
}
