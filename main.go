package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/raksit31667/example-go-api/config"
	"github.com/raksit31667/example-go-api/router"
)

func main() {
	e := echo.New()
	osGetter := &config.OsEnvGetter{}
	configProvider := config.ConfigProvider{Getter: osGetter}
	config := configProvider.GetConfig()
	
	db, err := sql.Open("postgres", config.Server.DBConnectionString)
	if err != nil {
		e.Logger.Fatal(err)
	}
	router.RegisterRoutes(e, db)
	address := fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {		
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()
	
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
