package handlers

import (
	"testing"

	"github.com/impr0ver/gophKeeper/internal/handlers/mocks"
	storageMocks "github.com/impr0ver/gophKeeper/internal/storage/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewServerHandlers(t *testing.T) {
	store := storageMocks.NewStorager(t)
	auth := mocks.NewAuthenticator(t)
	handlers := NewServerHandlers(store, auth)

	assert.NotEmpty(t, handlers)
}
