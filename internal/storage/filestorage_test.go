package storage

import (
	"testing"

	"github.com/impr0ver/gophKeeper/internal/serverconfig"

	"github.com/stretchr/testify/assert"
)

var filesPath = serverconfig.NewServerConfig().FilesStore

func TestNewFileStorage(t *testing.T) {
	assert.NotPanics(t, func() {
		newFileStorage(filesPath)
	})
}
