package handlers

import (
	"testing"

	"github.com/impr0ver/gophKeeper/internal/handlers/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewClientHandlers(t *testing.T) {
	conn := mocks.NewClientConnection(t)
	handlers := newClientHandlers(conn)
	assert.NotEmpty(t, handlers)
}
