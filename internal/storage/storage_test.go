package storage

import (
	"context"
	"testing"

	"github.com/impr0ver/gophKeeper/internal/storage/mocks"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewStorage(t *testing.T) {
	db, file := mocks.NewStorager(t), mocks.NewFileStorager(t)
	storage := NewStorage(db, file)

	assert.NotEmpty(t, storage)
}

func TestStorage_CreateUser(t *testing.T) {
	db, file := mocks.NewStorager(t), mocks.NewFileStorager(t)
	storage := NewStorage(db, file)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Create user",
			func() {
				db.On(
					"CreateUser",
					mock.AnythingOfType("userdata.UserCredentials"),
				).Return(nil)
			},
			func() {
				_ = storage.CreateUser(userdata.UserCredentials{
					Login:    "login",
					Password: "password",
				})
				db.AssertExpectations(t)
			},
		},
	}
	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}

func TestStorage_GetRecordsInfo(t *testing.T) {
	db, file := mocks.NewStorager(t), mocks.NewFileStorager(t)
	storage := NewStorage(db, file)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Get all records info",
			func() {
				db.On("GetRecordsInfo", context.Background()).Return([]userdata.Record{}, nil)
			},
			func() {
				_, _ = storage.GetRecordsInfo(context.Background())
				db.AssertExpectations(t)
			},
		},
	}
	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}

func TestStorage_LoginUser(t *testing.T) {
	db, file := mocks.NewStorager(t), mocks.NewFileStorager(t)
	storage := NewStorage(db, file)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Login user",
			func() {
				db.On(
					"LoginUser",
					mock.AnythingOfType("userdata.UserCredentials"),
				).Return(userdata.UserID(""), nil)
			},
			func() {
				_, _ = storage.LoginUser(userdata.UserCredentials{
					Login:    "login",
					Password: "password",
				})
				db.AssertExpectations(t)
				file.AssertExpectations(t)
			},
		},
	}
	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}

func TestStorage_CreateRecord(t *testing.T) {
	db, file := mocks.NewStorager(t), mocks.NewFileStorager(t)
	storage := NewStorage(db, file)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Create text record",
			func() {
				db.On(
					"CreateRecord",
					context.Background(),
					mock.AnythingOfType("userdata.Record"),
				).Return("", nil)
			},
			func() {
				_, _ = storage.CreateRecord(context.Background(), userdata.Record{
					ID:       "",
					Metadata: "",
					Type:     userdata.TypeText,
					Data:     nil,
				})
				db.AssertExpectations(t)
				file.AssertExpectations(t)
			},
		},
		{
			"Create file record",
			func() {
				db.On(
					"CreateRecord",
					context.Background(),
					mock.AnythingOfType("userdata.Record"),
				).Return("", nil)
				file.On(
					"CreateRecord",
					context.Background(),
					mock.AnythingOfType("userdata.Record"),
				).Return("", nil)
			},
			func() {
				_, _ = storage.CreateRecord(context.Background(), userdata.Record{
					ID:       "",
					Metadata: "",
					Type:     userdata.TypeFile,
					Data:     nil,
				})
				db.AssertExpectations(t)
				file.AssertExpectations(t)
			},
		},
	}
	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}

func TestStorage_GetRecord(t *testing.T) {
	db, file := mocks.NewStorager(t), mocks.NewFileStorager(t)
	storage := NewStorage(db, file)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Get file record",
			func() {
				db.On(
					"GetRecord",
					context.Background(),
					"",
				).Return(userdata.Record{Type: userdata.TypeFile}, nil)
				file.On(
					"GetRecord",
					mock.AnythingOfType("*context.valueCtx"),
					"",
				).Return(userdata.Record{}, nil)
			},
			func() {
				_, _ = storage.GetRecord(context.Background(), "")
				db.AssertExpectations(t)
				file.AssertExpectations(t)
			},
		},
		{
			"Get text record",
			func() {
				db.On("GetRecord", context.Background(), "").Return(userdata.Record{}, nil)
			},
			func() {
				_, _ = storage.GetRecord(context.Background(), "")
			},
		},
	}
	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}

func TestStorage_DeleteRecord(t *testing.T) {
	db, file := mocks.NewStorager(t), mocks.NewFileStorager(t)
	storage := NewStorage(db, file)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Delete file record",
			func() {
				db.On("DeleteRecord", context.Background(), "").Return(nil)
				file.On("DeleteRecord", context.Background(), "").Return(nil)
			},
			func() {
				_ = storage.DeleteRecord(context.Background(), "")
				db.AssertExpectations(t)
				file.AssertExpectations(t)
			},
		},
		{
			"Delete text record",
			func() {
				db.On("DeleteRecord", context.Background(), "").Return(userdata.Record{}, nil)
			},
			func() {
				_ = storage.DeleteRecord(context.Background(), "")
			},
		},
	}
	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}
