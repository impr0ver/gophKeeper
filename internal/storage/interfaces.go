package storage

import (
	"context"

	"github.com/impr0ver/gophKeeper/internal/userdata"
)

type DataBaseStorager interface {
	MigrateUP()
	CreateUser(credentials userdata.UserCredentials) error
	LoginUser(credentials userdata.UserCredentials) (userdata.UserID, error)
	GetRecordsInfo(ctx context.Context) ([]userdata.Record, error)
	CreateRecord(ctx context.Context, record userdata.Record) (string, error)
	GetRecord(ctx context.Context, recordID string) (userdata.Record, error)
	DeleteRecord(ctx context.Context, recordID string) error
}

// NewDBStorage connects to DB (interface).
func NewDBStorage(connectionURL string, migrateURL string) DataBaseStorager {
	return newDBStorage(connectionURL, migrateURL)
}

// FileStorager interface for storage, which can storage files.
//
//go:generate mockery --name FileStorager
type FileStorager interface {
	GetRecord(ctx context.Context, recordID string) (userdata.Record, error)
	CreateRecord(ctx context.Context, record userdata.Record) (string, error)
	DeleteRecord(ctx context.Context, recordID string) error
}

// NewFileStorage returns new file storage (interface).
func NewFileStorage(directory string) FileStorager {
	return newFileStorage(directory)
}

// Storager interface for storage, which can storage only text data.
//
//go:generate mockery --name Storager
type Storager interface {
	CreateUser(credentials userdata.UserCredentials) error
	LoginUser(credentials userdata.UserCredentials) (userdata.UserID, error)
	GetRecordsInfo(ctx context.Context) ([]userdata.Record, error)
	FileStorager
}
