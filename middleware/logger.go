package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	loggerContextKey = "logger"
	parentIDLogField = "parent-id"
	spanIDLogField   = "span-id"
)

func GetLogger(c echo.Context) *zap.Logger {
	return c.Get(loggerContextKey).(*zap.Logger)
}

func LogMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return log(next, logger)
	}
}

func log(next echo.HandlerFunc, logger *zap.Logger) func(c echo.Context) error {
	return func(c echo.Context) error {
		c.Set(loggerContextKey, loggerWithParentAndSpanID(c, logger))
		return next(c)
	}
}

func loggerWithParentAndSpanID(c echo.Context, logger *zap.Logger) *zap.Logger {
	parentID := c.Request().Header.Get("X-Request-ID")
	if parentID == "" {
		parentID = uuid.New().String()
	}
	spanID := uuid.New().String()
	return logger.With(zap.String(parentIDLogField, parentID), zap.String(spanIDLogField, spanID))
}
