package handlers

import (
	"context"

	"github.com/impr0ver/gophKeeper/internal/crypt"
	"github.com/impr0ver/gophKeeper/internal/storage"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	log "github.com/sirupsen/logrus"
)

// server struct for server handlers.
type server struct {
	Storage       storage.Storager
	Authenticator Authenticator
}

// newServerHandlers returns server handlers based on storage and authenticator (interface).
func newServerHandlers(s storage.Storager, a Authenticator) *server {
	return &server{
		Storage:       s,
		Authenticator: a,
	}
}

// LoginUser logins user by login and password.
func (s *server) LoginUser(credentials userdata.UserCredentials) (userdata.AuthToken, error) {
	if credentials.Login == "" || credentials.Password == "" {
		return "", ErrEmptyField
	}

	credentials.Password = crypt.PasswordHash(credentials)

	userID, err := s.Storage.LoginUser(credentials)
	if err != nil {
		log.Warnf("%s :: %v", "get user login fault", err)

		return "", err
	}

	authToken, errCreateToken := s.Authenticator.CreateToken(userID)
	if errCreateToken != nil {
		log.Warnf("%s :: %v", "create token fault", err)

		return "", storage.ErrUnknown
	}

	return authToken, nil
}

// CreateUser creates new user by login and password.
func (s *server) CreateUser(credentials userdata.UserCredentials) (userdata.AuthToken, error) {
	if credentials.Login == "" || credentials.Password == "" {
		return "", ErrEmptyField
	}

	if err := s.Storage.CreateUser(userdata.UserCredentials{
		Login:    credentials.Login,
		Password: crypt.PasswordHash(credentials),
	}); err != nil {
		log.Warnf("%s :: %v", "create new user fault", err)

		return "", err
	}

	return s.LoginUser(credentials)
}

// CreateRecord added record to storage.
func (s *server) CreateRecord(ctx context.Context, record userdata.Record) error {
	_, err := s.Storage.CreateRecord(ctx, record)
	return err
}

// GetRecordsInfo gets all records from storage.
func (s *server) GetRecordsInfo(ctx context.Context) ([]userdata.Record, error) {
	return s.Storage.GetRecordsInfo(ctx)
}

// GetRecord get record from storage by ID.
func (s *server) GetRecord(ctx context.Context, recordID string) (userdata.Record, error) {
	return s.Storage.GetRecord(ctx, recordID)
}

// DeleteRecord deletes record from storage.
func (s *server) DeleteRecord(ctx context.Context, recordID string) error {
	return s.Storage.DeleteRecord(ctx, recordID)
}
