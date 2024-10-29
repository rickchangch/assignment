package main

import (
	"assignment-pe/internal/config"
	"assignment-pe/internal/log"
	"assignment-pe/internal/postgre"
	"assignment-pe/internal/redis"
	"assignment-pe/internal/repo"
	"assignment-pe/internal/rest/controller"
	"assignment-pe/internal/rest/middleware"
	"assignment-pe/internal/rest/route"
	"assignment-pe/internal/rest/server"
	"assignment-pe/internal/service"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// nolint:funlen
func main() {
	config, err := config.Load()
	if err != nil {
		panic("load config failed: " + err.Error())
	}

	logger, err := log.NewLogger(log.AppLoggerConfig{
		Level:  config.Log.Level,
		Fields: log.Fields{"app": config.Server.Name},
	})
	if err != nil {
		panic("init logger failed: " + err.Error())
	}

	pgdb, err := postgre.NewPostgresDB(postgre.PostgresDBConfig{
		Username:    config.Connection.PostgreSQL.User,
		Password:    config.Connection.PostgreSQL.Password,
		Host:        config.Connection.PostgreSQL.Host,
		Database:    config.Connection.PostgreSQL.DB,
		MaxConn:     config.Connection.PostgreSQL.MaxConn,
		MaxConnIdle: config.Connection.PostgreSQL.MaxConnIdle,
	})
	if err != nil {
		panic("init postgresql failed: " + err.Error())
	}
	defer func() {
		if err := pgdb.Close(); err != nil {
			logger.WithError(err).Error("close postgresql failed")
		} else {
			logger.Info("postgresql closed")
		}
	}()

	redis, err := redis.NewRedis(redis.RedisConfig{
		Addr:         config.Connection.Redis.Host,
		Password:     config.Connection.Redis.Password,
		WriteTimeout: config.Connection.Redis.WriteTimeoutSec,
		ReadTimeout:  config.Connection.Redis.ReadTimeoutSec,
	})
	if err != nil {
		panic("init redis failed: " + err.Error())
	}
	defer func() {
		if err := redis.Close(); err != nil {
			logger.WithError(err).Error("close redis failed")
		} else {
			logger.Info("redis closed")
		}
	}()

	// Repo
	campaignRepo := repo.NewCampaignRepo()
	userRepo := repo.NewUserRepo()
	userCampaignRepo := repo.NewUserCampaignRepo()
	pointHistoryRepo := repo.NewPointHistoryRepo()

	// Service
	campaignSrv := service.NewCampaignService(redis, campaignRepo, userCampaignRepo)
	userSrv := service.NewUserService(userRepo)
	userCampaignSrv := service.NewUserCampaignService(userCampaignRepo)
	pointHistorySrv := service.NewPointHistoryService(pointHistoryRepo)
	swapSrv := service.NewSwapService(redis, campaignSrv, userCampaignSrv, pointHistorySrv)

	// Controller
	campaignCtrl := controller.NewCampaignController(campaignSrv)
	userCampaignCtrl := controller.NewUserCampaignController(userCampaignSrv, userSrv)
	pointHistoryCtrl := controller.NewPointHistoryController(pointHistorySrv)
	swapCtrl := controller.NewSwapController(swapSrv)

	// Server
	ginServer, engine := server.NewGinServer(config.Server.Port)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := ginServer.Shutdown(ctx); err != nil {
			logger.WithError(err).Error("shutdown gin server failed")
		} else {
			logger.Info("gin server closed")
		}
	}()

	middleware := middleware.NewMiddleware(logger, pgdb)
	route := route.NewRoute(
		engine, middleware,
		campaignCtrl, userCampaignCtrl, pointHistoryCtrl, swapCtrl,
	)
	route.Index()

	go func() {
		logger.WithField("port", config.Server.Port).Info("gin server start")
		if err := ginServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Error("unexpected stopped")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	if sig, ok := <-sigs; ok {
		logger.Infof("received signal: %v", sig)
	}
	logger.Info("start to shutdown")
	close(sigs)
}
