package main

import (
	"context"
	"fmt"
	"go-clean-arch/config"
	"go-clean-arch/helper-libs/loghelper"
	v1 "go-clean-arch/internal/api/v1"
	"go-clean-arch/internal/diregistry"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"

	"github.com/labstack/echo/v4"
	_ "github.com/swaggo/echo-swagger"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

// @title           Swagger Card Integrate API
// @version         1.0
// @description     This is a sample server card integrate proxy server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    hdbank.com.vn
// @contact.email  hoanglh7@hdbank.com.vn

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8280
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Getting configuration base on environment
	diregistry.BuildDIContainer()
	cfg := diregistry.GetDependency(diregistry.ConfigDIName).(*config.Config)

	err := loghelper.InitZap(cfg.App, cfg.Env, cfg.SensitiveFields)
	if err != nil {
		loghelper.Logger.Panic("Can't init zap logger", zap.Error(err))
	}

	loghelper.Logger.Infof("config: %+v", *cfg)

	httpServer := echo.New()
	if cfg.Env == "dev" {
		// use echoSwagger middleware to serve the API docs
		httpServer.GET("/swagger/*", echoSwagger.WrapHandler)
		// httpServer.Static("/swaggerui", "swaggerui")
	}
	httpServer.Use(echoprometheus.NewMiddleware("myapp"))   // adds middleware to gather metrics
	httpServer.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	httpServer.Server.Addr = fmt.Sprintf(":%d", cfg.HttpAddress)
	httpServer.Server.ReadTimeout = 10 * time.Minute
	httpServer.Server.WriteTimeout = 5 * time.Minute

	// Init route
	APIServer := diregistry.GetDependency(diregistry.ApiServerV1DIName).(v1.APIServer)
	v1publicRouter := httpServer.Group("/v1")
	APIServer.ConfigRoute(v1publicRouter)

	// Start the HTTP server
	go func() {
		if err := httpServer.StartServer(httpServer.Server); err != nil {
			if err == http.ErrServerClosed {
				httpServer.Logger.Info("shutting down the server")
			} else {
				httpServer.Logger.Errorf("error shutting down the server: ", err)
			}
		}
	}()
	loghelper.Logger.Infof("Gateway is started on port %d", cfg.HttpAddress)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	signals := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)
	// Grateful shutdown
	go func() {
		<-signals
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Stop HTTP Server
		if err := httpServer.Shutdown(ctx); err != nil {
			loghelper.Logger.WithContext(ctx).Errorf("Failed to shutdown gateway - err", err)
		}

		loghelper.Logger.Info("*****GRACEFUL SHUTTING DOWN*****")
		switch cfg.Env {
		case "prd":
			time.Sleep(15 * time.Second)
		case "dev":
			time.Sleep(3 * time.Second)
		}
		shutdown <- true
	}()
	<-shutdown
	loghelper.Logger.Info("*****SHUTDOWN*****")
}
