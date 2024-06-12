package handlers

import (
	"context"
	"github.com/impr0ver/gophKeeper/internal/jwtauth"
	"github.com/impr0ver/gophKeeper/internal/storage"
	"github.com/impr0ver/gophKeeper/internal/userdata"
	"time"
)

// ClientHandlers interface for Client.
type ClientHandlers interface {
	Login(credentials userdata.UserCredentials) error
	Register(credentials userdata.UserCredentials) error
	GetRecordsInfo() ([]userdata.Record, error)
	GetRecord(recordID string) (userdata.Record, error)
	CreateRecord(record userdata.Record) error
	DeleteRecord(recordID string) error
	SetAESKey(newAESKey string) error
}

// NewClientHandlers returns new client handlers (interface).
func NewClientHandlers(conn ClientConnection) ClientHandlers {
	return newClientHandlers(conn)
}

// Authenticator is interface for user authenticating. Should can creates tokens, and gets userIDs from them.
//
//go:generate mockery --name Authenticator
type Authenticator interface {
	CreateToken(userID userdata.UserID) (userdata.AuthToken, error)
	ValidateToken(token userdata.AuthToken) (userdata.UserID, error)
}

// NewAuthenticatorJWT gets new authenticatorJWT (interface).
func NewAuthenticatorJWT(secretKey []byte, expirationTime time.Duration) Authenticator {
	return jwtauth.NewAuthenticatorJWT(secretKey, expirationTime)
}

// ServerHandlers interface for server handlers
//
//go:generate mockery --name ServerHandlers
type ServerHandlers interface {
	LoginUser(credentials userdata.UserCredentials) (userdata.AuthToken, error)
	CreateUser(credentials userdata.UserCredentials) (userdata.AuthToken, error)
	GetRecordsInfo(ctx context.Context) ([]userdata.Record, error)
	GetRecord(ctx context.Context, recordID string) (userdata.Record, error)
	CreateRecord(ctx context.Context, record userdata.Record) error
	DeleteRecord(ctx context.Context, recordID string) error
}

// NewServerHandlers returns server handlers based on storage and authenticator.
func NewServerHandlers(s storage.Storager, a Authenticator) ServerHandlers {
	return newServerHandlers(s, a)
}

// ClientConnection describes client connection.
//
//go:generate mockery --name ClientConnection
type ClientConnection interface {
	Login(credentials userdata.UserCredentials) (string, error)
	Register(credentials userdata.UserCredentials) (string, error)
	GetRecordsInfo(token userdata.AuthToken) ([]userdata.Record, error)
	GetRecord(token userdata.AuthToken, recordID string) (userdata.Record, error)
	DeleteRecord(token userdata.AuthToken, recordID string) error
	CreateRecord(token userdata.AuthToken, record userdata.Record) error
}

// NewClientConnection connects to server and returning connection (interface).
func NewClientConnection(serverAddress string, clientCert string) ClientConnection {
	return newClientConn(serverAddress, clientCert)
}
