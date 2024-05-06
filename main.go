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
	"github.com/raksit31667/example-go-api/middleware"
	"github.com/raksit31667/example-go-api/router"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.LogMiddleware(logger))
	osGetter := &config.OsEnvGetter{}
	configProvider := config.ConfigProvider{Getter: osGetter}
	config := configProvider.GetConfig()

	db, err := sql.Open("postgres", config.Server.DBConnectionString)
	if err != nil {
		logger.Fatal("failed to open database connection", zap.Error(err))
	}
	router.RegisterRoutes(e, db)
	address := fmt.Sprintf("%s:%d", config.Server.Hostname, config.Server.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal("failed to shutdown server", zap.Error(err))
	}
}
