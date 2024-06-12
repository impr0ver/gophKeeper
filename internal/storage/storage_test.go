package storage

import (
	"testing"

	"github.com/impr0ver/gophKeeper/internal/storage/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	db, file := mocks.NewStorager(t), mocks.NewFileStorager(t)
	storage := NewStorage(db, file)

	assert.NotEmpty(t, storage)
}
