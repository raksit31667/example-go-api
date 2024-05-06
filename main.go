package main

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/raksit31667/example-go-api/config"
	"github.com/raksit31667/example-go-api/router"
	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()
	router.RegisterRoutes(e)
	osGetter := &config.OsEnvGetter{}
	configProvider := config.ConfigProvider{Getter: osGetter}
	config := configProvider.GetConfig()

	_, err := sql.Open("postgres", config.Server.DBConnectionString)
	if err != nil {
		e.Logger.Fatal(err)
	}

	address := fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port)
	e.Logger.Fatal(e.Start(address))
}
