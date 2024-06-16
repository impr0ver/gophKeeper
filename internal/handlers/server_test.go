package handlers

import (
	"context"
	"testing"

	"github.com/impr0ver/gophKeeper/internal/handlers/mocks"
	storMocks "github.com/impr0ver/gophKeeper/internal/storage/mocks"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

func TestNewServerHandlers(t *testing.T) {
	store := storMocks.NewStorager(t)
	auth := mocks.NewAuthenticator(t)
	handlers := NewServerHandlers(store, auth)

	assert.NotEmpty(t, handlers)
}

func TestServer_CreateUser(t *testing.T) {
	store := storMocks.NewStorager(t)
	auth := mocks.NewAuthenticator(t)
	handlers := NewServerHandlers(store, auth)

	tc := []struct {
		name string
		mock func()
		arg  userdata.UserCredentials
		want error
	}{
		{
			"User with creds",
			func() {
				store.On("CreateUser", userdata.UserCredentials{
					Login:    "Admin",
					Password: "b07e019b4662035489e1664afa63e28929a9df529f7a7fd6989e682e3cb695fd",
				}).Return(nil).Once()
				store.On("LoginUser", userdata.UserCredentials{
					Login:    "Admin",
					Password: "b07e019b4662035489e1664afa63e28929a9df529f7a7fd6989e682e3cb695fd",
				}).Return(userdata.UserID("userID"), nil).Once()
				auth.On("CreateToken", userdata.UserID("userID")).Return(userdata.AuthToken("token"), nil).Once()
			},
			userdata.UserCredentials{
				Login:    "Admin",
				Password: "password",
			},
			nil,
		},
		{
			"User with bad creds",
			func() {},
			userdata.UserCredentials{
				Login:    "",
				Password: "",
			},
			ErrEmptyField,
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		_, err := handlers.CreateUser(test.arg)
		assert.Equal(t, test.want, err)

		store.AssertExpectations(t)
		auth.AssertExpectations(t)
	}

}

func TestServer_LoginUser(t *testing.T) {
	store := storMocks.NewStorager(t)
	auth := mocks.NewAuthenticator(t)
	handlers := NewServerHandlers(store, auth)

	tc := []struct {
		name string
		mock func()
		arg  userdata.UserCredentials
		want error
	}{
		{
			"Login user with creds",
			func() {
				store.On("LoginUser", userdata.UserCredentials{
					Login:    "Admin",
					Password: "b07e019b4662035489e1664afa63e28929a9df529f7a7fd6989e682e3cb695fd",
				}).Return(userdata.UserID("userID"), nil).Once()
				auth.On("CreateToken", userdata.UserID("userID")).Return(userdata.AuthToken("token"), nil).Once()
			},
			userdata.UserCredentials{
				Login:    "Admin",
				Password: "password",
			},
			nil,
		},
		{
			"Login user with bad creds",
			func() {},
			userdata.UserCredentials{
				Login:    "",
				Password: "",
			},
			ErrEmptyField,
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		_, err := handlers.LoginUser(test.arg)
		assert.Equal(t, test.want, err)

		store.AssertExpectations(t)
		auth.AssertExpectations(t)
	}

}

func TestServer_GetRecordsInfo(t *testing.T) {
	store := storMocks.NewStorager(t)
	auth := mocks.NewAuthenticator(t)
	handlers := NewServerHandlers(store, auth)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Get all records",
			func() {
				store.On("GetRecordsInfo", mock.AnythingOfType("*context.valueCtx")).Return([]userdata.Record{}, nil).Once()

			},
			func() {
				md := metadata.Pairs("authToken", string("token"))
				ctx := metadata.NewIncomingContext(context.Background(), md)

				recInfo, err := handlers.GetRecordsInfo(ctx)
				assert.NoError(t, err)
				assert.Equal(t, []userdata.Record{}, recInfo)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()

		store.AssertExpectations(t)
		auth.AssertExpectations(t)
	}

}

func TestServer_GetRecord(t *testing.T) {
	store := storMocks.NewStorager(t)
	auth := mocks.NewAuthenticator(t)
	handlers := NewServerHandlers(store, auth)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Get record",
			func() {
				store.On("GetRecord", mock.AnythingOfType("*context.valueCtx"), "recordID").Return(userdata.Record{}, nil).Once()
			},
			func() {
				md := metadata.Pairs("authToken", string("token"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
				_, err := handlers.GetRecord(ctx, "recordID")
				assert.NoError(t, err)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()

		store.AssertExpectations(t)
		auth.AssertExpectations(t)
	}

}

func TestServer_CreateRecord(t *testing.T) {
	store := storMocks.NewStorager(t)
	auth := mocks.NewAuthenticator(t)
	handlers := NewServerHandlers(store, auth)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Create record",
			func() {
				store.On("CreateRecord", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("userdata.Record")).Return("", nil).Once()
			},
			func() {
				md := metadata.Pairs("authToken", string("token"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
				err := handlers.CreateRecord(ctx, userdata.Record{})
				assert.NoError(t, err)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()

		store.AssertExpectations(t)
		auth.AssertExpectations(t)
	}

}

func TestServer_DeleteRecord(t *testing.T) {
	store := storMocks.NewStorager(t)
	auth := mocks.NewAuthenticator(t)
	handlers := NewServerHandlers(store, auth)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Delete record",
			func() {
				store.On("DeleteRecord", mock.AnythingOfType("*context.valueCtx"), "recordID").Return(nil).Once()
			},
			func() {
				md := metadata.Pairs("authToken", string("token"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
				err := handlers.DeleteRecord(ctx, "recordID")
				assert.NoError(t, err)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()

		store.AssertExpectations(t)
		auth.AssertExpectations(t)
	}

}
