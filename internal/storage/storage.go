package storage

import (
	"context"
	"errors"

	"github.com/impr0ver/gophKeeper/internal/userdata"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

// Storage struct which saves to DB and file storage.
type Storage struct {
	DBStorage   Storager
	FileStorage FileStorager
}

// NewStorage returns new storage.
func NewStorage(DBStorage Storager, fileStorage FileStorager) *Storage {
	return &Storage{
		DBStorage:   DBStorage,
		FileStorage: fileStorage,
	}
}

// LoginUser check user login using DB storage.
func (s *Storage) LoginUser(credentials userdata.UserCredentials) (userdata.UserID, error) {
	return s.DBStorage.LoginUser(credentials)
}

// CreateUser creates new user and saves to DB storage.
func (s *Storage) CreateUser(credentials userdata.UserCredentials) error {
	return s.DBStorage.CreateUser(credentials)
}

// GetRecordsInfo gets all records from user from DB storage.
func (s *Storage) GetRecordsInfo(ctx context.Context) ([]userdata.Record, error) {
	return s.DBStorage.GetRecordsInfo(ctx)
}

// CreateRecord creates record, saves to DB. If record type is file, saves to file storage too.
func (s *Storage) CreateRecord(ctx context.Context, record userdata.Record) (string, error) {
	data := record.Data

	if record.Type == userdata.TypeFile {
		record.Data = nil
	}

	id, err := s.DBStorage.CreateRecord(ctx, record)
	if err != nil {
		log.Infoln(err)

		return "", err
	}

	if record.Type == userdata.TypeFile {
		record.ID = id
		record.Data = data
		_, err = s.FileStorage.CreateRecord(ctx, record)
		return "", err
	}

	return id, nil
}

// DeleteRecord deletes record from DB storage. If record type is file, delete file from storage.
func (s *Storage) DeleteRecord(ctx context.Context, recordID string) error {
	err := s.DBStorage.DeleteRecord(ctx, recordID)
	if err != nil {
		log.Infoln(err)

		return err
	}

	err = s.FileStorage.DeleteRecord(ctx, recordID)
	if !errors.Is(err, ErrNotFound) && err != nil {
		return ErrUnknown
	}

	return nil
}

// GetRecord gets record from DB or file storage.
func (s *Storage) GetRecord(ctx context.Context, recordID string) (userdata.Record, error) {
	record, err := s.DBStorage.GetRecord(ctx, recordID)
	if err != nil {
		log.Infoln(err)

		return record, err
	}

	if record.Type == userdata.TypeFile {

		md := metadata.Pairs(
			"recordMetadata", record.Metadata,
		)
		ctx := metadata.NewIncomingContext(ctx, md)
		return s.FileStorage.GetRecord(ctx, recordID)
	}

	return record, nil
}
