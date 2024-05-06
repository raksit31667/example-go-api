package router

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/raksit31667/example-go-api/user"
)

func RegisterRoutes(e *echo.Echo, db *sql.DB) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	userHandler := user.NewHandler(db)

	e.POST("/users", func(c echo.Context) error {
		return userHandler.Create(c)
	})
}
