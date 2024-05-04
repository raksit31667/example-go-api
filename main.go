package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/raksit31667/example-go-api/config"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	osGetter := &config.OsEnvGetter{}
	configProvider := config.ConfigProvider{Getter: osGetter}
	config := configProvider.GetConfig()
	address := fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port)
	e.Logger.Fatal(e.Start(address))
}
