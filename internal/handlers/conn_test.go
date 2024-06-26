package handlers

import (
	"context"
	"testing"

	"github.com/impr0ver/gophKeeper/internal/clientconfig"
	"github.com/impr0ver/gophKeeper/internal/handlers/mocks"
	"github.com/impr0ver/gophKeeper/internal/serverconfig"
	"github.com/impr0ver/gophKeeper/internal/storage"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	var serverCfg = serverconfig.ServerConfig{}
	serverCfg.ServerCert = "../../cmd/cert/server-cert.pem"
	serverCfg.ServerKey = "../../cmd/cert/server-key.pem"
	serverCfg.ListenAddr = "127.0.0.1:9000"

	var clientCfg = clientconfig.ClientConfig{}
	clientCfg.ServerAddress = "127.0.0.1:9000"
	clientCfg.ClientCert = "../../cmd/cert/ca-cert.pem"

	auth := mocks.NewAuthenticator(t)
	handlers := mocks.NewServerHandlers(t)

	server := NewServerConn(handlers, auth, serverCfg.ServerCert, serverCfg.ServerKey, serverCfg.ServerConsoleLog)
	
	ctx, cancel := context.WithCancel(context.Background())
	server.Start(ctx, serverCfg.ListenAddr)
	
	client := newClientConn(clientCfg.ServerAddress, clientCfg.ClientCert)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Create user",
			func() {
				handlers.On("CreateUser", userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				}).Return(userdata.AuthToken("token"), nil).Once()
			},
			func() {
				token, err := client.Register(userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				})
				assert.NoError(t, err)
				assert.Equal(t, "token", token)
			},
		},
		{
			"Create user, but error",
			func() {
				handlers.On("CreateUser", userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				}).Return(userdata.AuthToken("token"), storage.ErrLoginExists).Once()
			},
			func() {
				token, err := client.Register(userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				})
				assert.Equal(t, storage.ErrLoginExists, err)
				assert.Empty(t, token)
			},
		},
		{
			"Create user, but unknown error",
			func() {
				handlers.On("CreateUser", userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				}).Return(userdata.AuthToken("token"), storage.ErrUnknown).Once()
			},
			func() {
				token, err := client.Register(userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				})
				assert.Equal(t, storage.ErrUnknown, err)
				assert.Empty(t, token)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		handlers.AssertExpectations(t)
	}

	cancel()
	server.Stop()
	
}

func TestLoginUser(t *testing.T) {
	var serverCfg = serverconfig.ServerConfig{}
	serverCfg.ServerCert = "../../cmd/cert/server-cert.pem"
	serverCfg.ServerKey = "../../cmd/cert/server-key.pem"
	serverCfg.ListenAddr = "127.0.0.1:9000"

	var clientCfg = clientconfig.ClientConfig{}
	clientCfg.ServerAddress = "127.0.0.1:9000"
	clientCfg.ClientCert = "../../cmd/cert/ca-cert.pem"

	auth := mocks.NewAuthenticator(t)
	client := newClientConn(clientCfg.ServerAddress, clientCfg.ClientCert)
	handlers := mocks.NewServerHandlers(t)

	server := NewServerConn(handlers, auth, serverCfg.ServerCert, serverCfg.ServerKey, serverCfg.ServerConsoleLog)
	ctx, cancel := context.WithCancel(context.Background())
	server.Start(ctx, serverCfg.ListenAddr)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Login user",
			func() {
				handlers.On("LoginUser", userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				}).Return(userdata.AuthToken("token"), nil).Once()
			},
			func() {
				token, err := client.Login(userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				})
				assert.NoError(t, err)
				assert.Equal(t, "token", token)
			},
		},
		{
			"Login user, but error",
			func() {
				handlers.On("LoginUser", userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				}).Return(userdata.AuthToken("token"), storage.ErrWrongCredentials).Once()
			},
			func() {
				token, err := client.Login(userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				})
				assert.Equal(t, storage.ErrWrongCredentials, err)
				assert.Empty(t, token)
			},
		},
		{
			"Create user, but unknown error",
			func() {
				handlers.On("LoginUser", userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				}).Return(userdata.AuthToken("token"), storage.ErrUnknown).Once()
			},
			func() {
				token, err := client.Login(userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				})
				assert.Equal(t, storage.ErrUnknown, err)
				assert.Empty(t, token)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		handlers.AssertExpectations(t)
	}

	cancel()
	server.Stop()
}

func TestGetRecordsInfo(t *testing.T) {
	var serverCfg = serverconfig.ServerConfig{}
	serverCfg.ServerCert = "../../cmd/cert/server-cert.pem"
	serverCfg.ServerKey = "../../cmd/cert/server-key.pem"
	serverCfg.ListenAddr = "127.0.0.1:9000"

	var clientCfg = clientconfig.ClientConfig{}
	clientCfg.ServerAddress = "127.0.0.1:9000"
	clientCfg.ClientCert = "../../cmd/cert/ca-cert.pem"

	auth := mocks.NewAuthenticator(t)
	client := newClientConn(clientCfg.ServerAddress, clientCfg.ClientCert)
	handlers := mocks.NewServerHandlers(t)

	server := NewServerConn(handlers, auth, serverCfg.ServerCert, serverCfg.ServerKey, serverCfg.ServerConsoleLog)
	ctx, cancel := context.WithCancel(context.Background())
	server.Start(ctx, serverCfg.ListenAddr)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Get all records",
			func() {
				handlers.On("GetRecordsInfo", mock.AnythingOfType("*context.valueCtx")).
					Return([]userdata.Record{}, nil).Once()
				auth.On("ValidateToken", userdata.AuthToken("token")).Return(userdata.UserID("userID"), nil).Once()
			},
			func() {
				_, err := client.GetRecordsInfo("token")
				assert.NoError(t, err)
			},
		},
		{
			"Get all records, but error",
			func() {
				handlers.On("GetRecordsInfo", mock.AnythingOfType("*context.valueCtx")).
					Return([]userdata.Record{}, storage.ErrUnauthenticated).Once()
				auth.On("ValidateToken", userdata.AuthToken("token")).Return(userdata.UserID("userID"), nil).Once()
			},
			func() {
				_, err := client.GetRecordsInfo("token")
				assert.Equal(t, storage.ErrUnauthenticated, err)
			},
		},
		{
			"Get all records, but unknown error",
			func() {
				handlers.On("GetRecordsInfo", mock.AnythingOfType("*context.valueCtx")).
					Return([]userdata.Record{}, storage.ErrUnknown).Once()
				auth.On("ValidateToken", userdata.AuthToken("token")).Return(userdata.UserID("userID"), nil).Once()
			},
			func() {
				_, err := client.GetRecordsInfo("token")
				assert.Equal(t, storage.ErrUnknown, err)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		handlers.AssertExpectations(t)
	}

	cancel()
	server.Stop()
}

func TestGetRecord(t *testing.T) {
	var serverCfg = serverconfig.ServerConfig{}
	serverCfg.ServerCert = "../../cmd/cert/server-cert.pem"
	serverCfg.ServerKey = "../../cmd/cert/server-key.pem"
	serverCfg.ListenAddr = "127.0.0.1:9000"

	var clientCfg = clientconfig.ClientConfig{}
	clientCfg.ServerAddress = "127.0.0.1:9000"
	clientCfg.ClientCert = "../../cmd/cert/ca-cert.pem"

	auth := mocks.NewAuthenticator(t)
	client := newClientConn(clientCfg.ServerAddress, clientCfg.ClientCert)
	handlers := mocks.NewServerHandlers(t)

	server := NewServerConn(handlers, auth, serverCfg.ServerCert, serverCfg.ServerKey, serverCfg.ServerConsoleLog)
	ctx, cancel := context.WithCancel(context.Background())
	server.Start(ctx, serverCfg.ListenAddr)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Get record",
			func() {
				handlers.On("GetRecord", mock.AnythingOfType("*context.valueCtx"), "recordID").
					Return(userdata.Record{}, nil).Once()
				auth.On("ValidateToken", userdata.AuthToken("token")).Return(userdata.UserID("userID"), nil).Once()
			},
			func() {
				_, err := client.GetRecord("token", "recordID")
				assert.NoError(t, err)
			},
		},
		{
			"Get record, but wrong ID.",
			func() {
				handlers.On("GetRecord", mock.AnythingOfType("*context.valueCtx"), "recordID").
					Return(userdata.Record{}, storage.ErrNotFound).Once()
				auth.On("ValidateToken", userdata.AuthToken("token")).Return(userdata.UserID("userID"), nil).Once()
			},
			func() {
				_, err := client.GetRecord("token", "recordID")
				assert.Equal(t, storage.ErrNotFound, err)
			},
		},
		{
			"Get record, but unknown error.",
			func() {
				handlers.On("GetRecord", mock.AnythingOfType("*context.valueCtx"), "recordID").
					Return(userdata.Record{}, storage.ErrUnknown).Once()
				auth.On("ValidateToken", userdata.AuthToken("token")).Return(userdata.UserID("userID"), nil).Once()
			},
			func() {
				_, err := client.GetRecord("token", "recordID")
				assert.Equal(t, storage.ErrUnknown, err)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		handlers.AssertExpectations(t)
	}

	cancel()
	server.Stop()
}

func TestCreateRecord(t *testing.T) {
	var serverCfg = serverconfig.ServerConfig{}
	serverCfg.ServerCert = "../../cmd/cert/server-cert.pem"
	serverCfg.ServerKey = "../../cmd/cert/server-key.pem"
	serverCfg.ListenAddr = "127.0.0.1:9000"

	var clientCfg = clientconfig.ClientConfig{}
	clientCfg.ServerAddress = "127.0.0.1:9000"
	clientCfg.ClientCert = "../../cmd/cert/ca-cert.pem"

	auth := mocks.NewAuthenticator(t)
	client := newClientConn(clientCfg.ServerAddress, clientCfg.ClientCert)
	handlers := mocks.NewServerHandlers(t)

	server := NewServerConn(handlers, auth, serverCfg.ServerCert, serverCfg.ServerKey, serverCfg.ServerConsoleLog)
	ctx, cancel := context.WithCancel(context.Background())
	server.Start(ctx, serverCfg.ListenAddr)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Create record.",
			func() {
				handlers.On(
					"CreateRecord",
					mock.AnythingOfType("*context.valueCtx"),
					userdata.Record{},
				).Return(nil).Once()
				auth.On("ValidateToken", userdata.AuthToken("token")).Return(userdata.UserID("userID"), nil).Once()
			},
			func() {
				err := client.CreateRecord("token", userdata.Record{})
				assert.NoError(t, err)
			},
		},
		{
			"Create record, but unknown error.",
			func() {
				handlers.On(
					"CreateRecord",
					mock.AnythingOfType("*context.valueCtx"),
					userdata.Record{},
				).Return(storage.ErrUnknown).Once()
				auth.On("ValidateToken", userdata.AuthToken("token")).Return(userdata.UserID("userID"), nil).Once()
			},
			func() {
				err := client.CreateRecord("token", userdata.Record{})
				assert.Equal(t, storage.ErrUnknown, err)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		handlers.AssertExpectations(t)
	}

	cancel()
	server.Stop()
}

func TestDeleteRecord(t *testing.T) {
	var serverCfg = serverconfig.ServerConfig{}
	serverCfg.ServerCert = "../../cmd/cert/server-cert.pem"
	serverCfg.ServerKey = "../../cmd/cert/server-key.pem"
	serverCfg.ListenAddr = "127.0.0.1:9000"

	var clientCfg = clientconfig.ClientConfig{}
	clientCfg.ServerAddress = "127.0.0.1:9000"
	clientCfg.ClientCert = "../../cmd/cert/ca-cert.pem"

	auth := mocks.NewAuthenticator(t)
	client := newClientConn(clientCfg.ServerAddress, clientCfg.ClientCert)
	handlers := mocks.NewServerHandlers(t)

	server := NewServerConn(handlers, auth, serverCfg.ServerCert, serverCfg.ServerKey, serverCfg.ServerConsoleLog)
	ctx, cancel := context.WithCancel(context.Background())
	server.Start(ctx, serverCfg.ListenAddr)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Delete record.",
			func() {
				handlers.On(
					"DeleteRecord",
					mock.AnythingOfType("*context.valueCtx"),
					"recordID",
				).Return(nil).Once()
				auth.On("ValidateToken", userdata.AuthToken("token")).Return(userdata.UserID("userID"), nil).Once()
			},
			func() {
				err := client.DeleteRecord("token", "recordID")
				assert.NoError(t, err)
			},
		},
		{
			"Delete record, but not found.",
			func() {
				handlers.On(
					"DeleteRecord",
					mock.AnythingOfType("*context.valueCtx"),
					"recordID",
				).Return(storage.ErrNotFound).Once()
				auth.On("ValidateToken", userdata.AuthToken("token")).Return(userdata.UserID("userID"), nil).Once()
			},
			func() {
				err := client.DeleteRecord("token", "recordID")
				assert.Equal(t, storage.ErrNotFound, err)
			},
		},
		{
			"Delete record, but unknown error.",
			func() {
				handlers.On(
					"DeleteRecord",
					mock.AnythingOfType("*context.valueCtx"),
					"recordID",
				).Return(storage.ErrUnknown).Once()
				auth.On("ValidateToken", userdata.AuthToken("token")).Return(userdata.UserID("userID"), nil).Once()
			},
			func() {
				err := client.DeleteRecord("token", "recordID")
				assert.Equal(t, storage.ErrUnknown, err)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		handlers.AssertExpectations(t)
	}

	cancel()
	server.Stop()
}
