package main

import (
	"github.com/impr0ver/gophKeeper/internal/clientconfig"
	"github.com/impr0ver/gophKeeper/internal/clientwork"
	"github.com/impr0ver/gophKeeper/internal/handlers"
	"github.com/impr0ver/gophKeeper/internal/logger"
	log "github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	cfg := clientconfig.NewClientConfig()

	logger.NewLogrusLogger()
	
	buildInfo()
	
	conn := handlers.NewClientConnection(cfg.ServerAddress, cfg.ClientCert)
	handlers := handlers.NewClientHandlers(conn)

	termUserInterface := clientwork.NewTUI(handlers, cfg.MaxFileSize)

	//TUI close app via Ctrl+C (via method "app.Stop()" in lib "tview" is not work properly)
	err := termUserInterface.Run()
	if err == nil {
		log.Info("Stop client (TUI) successfully")
	}
}

func buildInfo() {
	log.Infof("Build client version: %s", buildVersion)
	log.Infof("Build client date: %s", buildDate)
	log.Infof("Build client commit: %s", buildCommit)
}
