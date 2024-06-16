package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/impr0ver/gophKeeper/internal/handlers"
	"github.com/impr0ver/gophKeeper/internal/logger"
	"github.com/impr0ver/gophKeeper/internal/serverconfig"
	"github.com/impr0ver/gophKeeper/internal/storage"
	log "github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	cfg := serverconfig.NewServerConfig()
	
	logger.NewLogrusLogger()
	sLogger := logger.NewSugarLogger()
	
	buildInfo()
	
	dataBase := storage.NewDBStorage(cfg.DatabaseDSN, cfg.MigrationsURL)
	dataBase.MigrateUP()
	files := storage.NewFileStorage(cfg.FilesStore)
	
	stor := storage.NewStorage(dataBase, files)

	jwtAuth := handlers.NewAuthenticatorJWT([]byte(cfg.JWTAuth.SecretJWT), cfg.JWTAuth.ExpirationTime)
	h := handlers.NewServerHandlers(stor, jwtAuth)
	server := handlers.NewServerConn(h, jwtAuth, cfg.ServerCert, cfg.ServerKey, cfg.ServerConsoleLog)

	go server.Start(context.Background(), cfg.ListenAddr)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint

	server.Stop()
	sLogger.Info("gRPC server is gracefully stop!")
	log.Info("gRPC server is gracefully stop!")
}

func buildInfo() {
	sLogger := logger.NewSugarLogger()
	sLogger.Infof("Build server version: %s", buildVersion)
	sLogger.Infof("Build server date: %s", buildDate)
	sLogger.Infof("Build server commit: %s", buildCommit)

	log.Infof("Build server version: %s", buildVersion)
	log.Infof("Build server date: %s", buildDate)
	log.Infof("Build server commit: %s", buildCommit)
}
