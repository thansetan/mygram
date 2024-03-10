package main

import (
	"context"
	"final-project/helper"
	"final-project/lib/config"
	"final-project/lib/database"
	"final-project/lib/logging"
	"final-project/middleware"
	"final-project/routes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	docs "final-project/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title MyGram
// @version 1.0
// @description MyGram (submission for Hacktiv8-Golang final-project)
// @securityDefinitions.apikey BearerToken
// @in header
// @name Authorization
// @description Bearer token for authentication. Format: Bearer {token}
func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "json-config", "config.json", "path to json config file")
	flag.Parse()

	logger := logging.New(os.Stderr)
	conf, err := config.Load(configFilePath)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	helper.JWTSecret = []byte(conf.App.JWTSecret)
	helper.JWTExpiresIn = helper.GetJWTExpiresIn(conf.App.JWTExpiresIn, time.Hour)

	db, err := database.New(conf.DB)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	api := http.NewServeMux()

	{
		routes.InitUserRoutes(api, db, logger)
		routes.InitPhotoRoutes(api, db, logger)
		routes.InitLikeRoutes(api, db, logger)
		routes.InitCommentRoutes(api, db, logger)
		routes.InitSocialMediaRoutes(api, db, logger)
	}

	r := http.NewServeMux()
	docs.SwaggerInfo.BasePath = conf.App.BasePath
	{
		r.Handle(conf.App.BasePath, middleware.Logging(http.StripPrefix(strings.TrimSuffix(conf.App.BasePath, "/"), api)))
		r.HandleFunc("GET /swagger/", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))
	}

	server := new(http.Server)
	server.Addr = fmt.Sprintf("%s:%d", conf.App.Host, conf.App.Port)
	server.Handler = middleware.Recover(r)

	logger.Info("Starting server...", "addr", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	<-ctx.Done()

	logger.Info("Shutting down server...", "addr", server.Addr)
	err = server.Shutdown(context.Background())
	if err != nil {
		logger.Error(err.Error())
	}
}
