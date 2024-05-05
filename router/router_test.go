package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/labstack/echo/v4"
)

func TestRegisterRoutes(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	e.ServeHTTP(response, request)
	RegisterRoutes(e)
	routes := e.Routes()
	want := 1
	if len(routes) != 1 {
		t.Errorf("got %d routes but want %d routes", len(routes), want)
	}
}