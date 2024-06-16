package handlers

import (
	"testing"

	"github.com/impr0ver/gophKeeper/internal/handlers/mocks"
	"github.com/impr0ver/gophKeeper/internal/storage"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewClientHandlers(t *testing.T) {
	conn := mocks.NewClientConnection(t)
	handlers := newClientHandlers(conn)
	assert.NotEmpty(t, handlers)
}

func TestClient_Register(t *testing.T) {
	conn := mocks.NewClientConnection(t)
	handlers := newClientHandlers(conn)
	handlers.AESKey = "hello"

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Register with good credentials",
			func() {
				conn.On("Register", userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
					AESKey:   "hello",
				}).Return("token", nil).Once()
			},
			func() {
				err := handlers.Register(userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
					AESKey:   "hello",
				})
				assert.NoError(t, err)
				assert.Equal(t, userdata.AuthToken("token"), handlers.authToken)
				assert.Equal(
					t,
					"hello",
					handlers.AESKey,
				)
			},
		},
		{
			"Register with bad credentials",
			func() {},
			func() {
				handlers.authToken = ""
				err := handlers.Register(userdata.UserCredentials{
					Login:    "",
					Password: "",
					AESKey:   "hello",
				})
				assert.Equal(t, ErrEmptyField, err)
				assert.Empty(t, handlers.authToken)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		conn.AssertExpectations(t)
	}
}

func TestClient_Login(t *testing.T) {
	conn := mocks.NewClientConnection(t)
	handlers := newClientHandlers(conn)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Login with good credentials",
			func() {
				conn.On("Login", userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
					AESKey:   "hello",
				}).Return("token", nil).Once()
			},
			func() {
				err := handlers.Login(userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
					AESKey:   "hello",
				})
				assert.NoError(t, err)
				assert.Equal(t, userdata.AuthToken("token"), handlers.authToken)
				assert.Equal(
					t,
					"hello",
					handlers.AESKey,
				)
			},
		},
		{
			"Login with bad credentials",
			func() {},
			func() {
				handlers.authToken = ""
				err := handlers.Login(userdata.UserCredentials{
					Login:    "",
					Password: "",
					AESKey:   "hello",
				})
				assert.Equal(t, ErrEmptyField, err)
				assert.Empty(t, handlers.authToken)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		conn.AssertExpectations(t)
	}
}

func TestClient_GetRecordsInfo(t *testing.T) {
	conn := mocks.NewClientConnection(t)
	handlers := newClientHandlers(conn)
	handlers.authToken = "token"

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Get records info",
			func() {
				conn.On(
					"GetRecordsInfo",
					userdata.AuthToken("token")).Return([]userdata.Record{},
					nil,
				).Once()
			},
			func() {
				records, err := handlers.GetRecordsInfo()
				assert.NoError(t, err)
				assert.Equal(t, []userdata.Record{}, records)
			},
		},
		{
			"Get records info, but return error",
			func() {
				conn.On(
					"GetRecordsInfo",
					userdata.AuthToken("token")).Return([]userdata.Record{},
					storage.ErrUnauthenticated,
				).Once()
			},
			func() {
				records, err := handlers.GetRecordsInfo()
				assert.Equal(t, storage.ErrUnauthenticated, err)
				assert.Equal(t, []userdata.Record{}, records)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		conn.AssertExpectations(t)
	}
}

func TestClient_GetRecord(t *testing.T) {
	conn := mocks.NewClientConnection(t)
	handlers := newClientHandlers(conn)
	handlers.authToken = "token"
	handlers.AESKey = "hello"

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Get record, but not found",
			func() {
				conn.On("GetRecord", userdata.AuthToken("token"), "1").
					Return(userdata.Record{}, storage.ErrNotFound).Once()
			},
			func() {
				record, err := handlers.GetRecord("1")
				assert.Equal(t, storage.ErrNotFound, err)
				assert.Empty(t, record)
			},
		},
		{
			"Get record, but unknown error",
			func() {
				conn.On("GetRecord", userdata.AuthToken("token"), "1").
					Return(userdata.Record{}, storage.ErrUnknown).Once()
			},
			func() {
				record, err := handlers.GetRecord("1")
				assert.Equal(t, storage.ErrUnknown, err)
				assert.Empty(t, record)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		conn.AssertExpectations(t)
	}
}

func TestClient_DeleteRecord(t *testing.T) {
	conn := mocks.NewClientConnection(t)
	handlers := newClientHandlers(conn)
	handlers.authToken = "token"
	handlers.AESKey = "hello"

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Delete record",
			func() {
				conn.On("DeleteRecord", userdata.AuthToken("token"), "1").Return(nil).Once()
			},
			func() {
				err := handlers.DeleteRecord("1")
				assert.NoError(t, err)
			},
		},
		{
			"Delete record, but will return error",
			func() {
				conn.On("DeleteRecord", userdata.AuthToken("token"), "1").Return(storage.ErrNotFound).Once()
			},
			func() {
				err := handlers.DeleteRecord("1")
				assert.Equal(t, storage.ErrNotFound, err)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		conn.AssertExpectations(t)
	}
}

func TestClient_CreateRecord(t *testing.T) {
	conn := mocks.NewClientConnection(t)
	handlers := newClientHandlers(conn)
	handlers.authToken = "token"
	handlers.AESKey = "masterkey"

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Create record",
			func() {
				conn.On(
					"CreateRecord",
					userdata.AuthToken("token"),
					mock.AnythingOfType("userdata.Record"),
				).Return(nil).Once()
			},
			func() {
				err := handlers.CreateRecord(userdata.Record{
					Data: []byte("hello!"),
				})
				assert.NoError(t, err)
			},
		},
		{
			"Create record, but user not authored",
			func() {
				conn.On(
					"CreateRecord",
					userdata.AuthToken("token"),
					mock.AnythingOfType("userdata.Record"),
				).Return(storage.ErrUnauthenticated).Once()
			},
			func() {
				err := handlers.CreateRecord(userdata.Record{
					Data: []byte("hello!"),
				})
				assert.Equal(t, storage.ErrUnauthenticated, err)
			},
		},
		{
			"Create record, but return unknown error",
			func() {
				conn.On(
					"CreateRecord",
					userdata.AuthToken("token"),
					mock.AnythingOfType("userdata.Record"),
				).Return(storage.ErrUnknown).Once()
			},
			func() {
				err := handlers.CreateRecord(userdata.Record{
					Data: []byte("hello!"),
				})
				assert.Equal(t, storage.ErrUnknown, err)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		conn.AssertExpectations(t)
	}
}
