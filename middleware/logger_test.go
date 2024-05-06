package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogMiddleware(t *testing.T) {
	t.Run("should set logger with default parent and span id to context", func(t *testing.T) {
		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		logger := zap.New(observedZapCore)
		e := echo.New()
		e.Use(LogMiddleware(logger))

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		handler := logHelloWorldHandler()
		err := log(handler, logger)(c)

		assertNoError(t, err)
		assertLoggerInContext(t, c)
		got := getFieldsInLogContext(t, observedLogs)
		if got[parentIDLogField] == "" {
			t.Error("expected parent-id in log context but got none")
		}
		if got[spanIDLogField] == "" {
			t.Error("expected span-id in log context but got none")
		}
	})

	t.Run("should set parent id with X-Request-ID to context given X-Request-ID exists", func(t *testing.T) {
		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		logger := zap.New(observedZapCore)
		e := echo.New()
		e.Use(LogMiddleware(logger))

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		c := e.NewContext(request, response)
		c.Request().Header.Set("X-Request-ID", "request-id")
		handler := logHelloWorldHandler()
		err := log(handler, logger)(c)

		assertNoError(t, err)
		assertLoggerInContext(t, c)
		got := getFieldsInLogContext(t, observedLogs)
		if got[parentIDLogField] != "request-id" {
			t.Errorf("expected parent-id to be request-id but got %s", got[parentIDLogField])
		}
		if got[spanIDLogField] == "" {
			t.Error("expected span-id in log context but got none")
		}
	})
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func assertLoggerInContext(t testing.TB, c echo.Context) {
	t.Helper()
	if GetLogger(c) == nil {
		t.Error("expected logger in context but got none")
	}
}

func logHelloWorldHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := GetLogger(c)
		logger.Info("Hello, World!")
		return c.String(http.StatusOK, "Hello, World!")
	}
}

func getFieldsInLogContext(t testing.TB, observedLogs *observer.ObservedLogs) map[string]string {
	t.Helper()
	got := make(map[string]string)
	for _, f := range observedLogs.All()[0].Context {
		got[f.Key] = f.String
	}

	return got
}
