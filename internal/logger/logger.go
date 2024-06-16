package logger

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// NewSugarLogger init zap sugar logger
func NewSugarLogger() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	//SugaredLogger registrator
	sugar := *logger.Sugar()
	return &sugar
}

// NewLogrusLogger init logrus logger
func NewLogrusLogger() {
	fileName := time.Now().Format("2006-01-02")
	file, err := os.OpenFile(fmt.Sprintf("%s%s", fileName, ".log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	// Log as JSON format.
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	log.SetOutput(file)

	// Display filename and line number
	log.SetReportCaller(true)

	// Log level.
	log.SetLevel(log.InfoLevel)
}
