package logger

import (
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogrusLogger(t *testing.T) {
	NewLogrusLogger()

	log.Info("testtesttest")

	fileName := time.Now().Format("2006-01-02")
	fullName := fmt.Sprintf("%s%s", fileName, ".log")
	bData, err := os.ReadFile(fullName)
	assert.NoError(t, err)

	assert.NotEmpty(t, bData)

	os.Remove(fullName)
}

func TestNewSugarLogger(t *testing.T){
	sLogger := NewSugarLogger()
	assert.NotEmpty(t, sLogger)
}
