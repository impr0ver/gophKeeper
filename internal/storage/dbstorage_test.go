package storage

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/impr0ver/gophKeeper/internal/userdata"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

// Init the flags before start unittests
var _ = func() bool {
	testing.Init()
	return true
}()

func TestNewDBStorage(t *testing.T) {
	assert.NotPanics(t, func() {
		NewDBStorage("", "")
	})
}

func TestDBStorage_CreateUser(t *testing.T) {
	storage := newDBStorage("", "")

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	storage.DB = db

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Create user with good credentials (doesn't exists)",
			func() {
				mock.ExpectQuery(
					`SELECT COUNT(*) FROM users WHERE login = $1`,
				).WithArgs("my_login").WillReturnRows(sqlmock.NewRows(
					[]string{"count"}).AddRow(0),
				)
				mock.ExpectExec(
					`INSERT INTO users (login, password) VALUES ($1, $2)`,
				).WithArgs("my_login", "my_password").WillReturnResult(
					sqlmock.NewResult(0, 1),
				)
			},
			func() {
				err := storage.CreateUser(userdata.UserCredentials{
					Login:    "my_login",
					Password: "my_password",
				})
				assert.NoError(t, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Create user with good credentials (doesn't exists), but DB will return error",
			func() {
				mock.ExpectQuery(
					`SELECT COUNT(*) FROM users WHERE login = $1`,
				).WithArgs("my_login").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectExec(
					`INSERT INTO users (login, password) VALUES ($1, $2)`,
				).WithArgs("my_login", "my_password").WillReturnError(errors.New("some DB error"))
			},
			func() {
				err := storage.CreateUser(userdata.UserCredentials{
					Login:    "my_login",
					Password: "my_password",
				})
				assert.Equal(t, ErrUnknown, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Create user with good credentials (already exists)",
			func() {
				mock.ExpectQuery(
					`SELECT COUNT(*) FROM users WHERE login = $1`,
				).WithArgs("my_login").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			func() {
				err := storage.CreateUser(userdata.UserCredentials{
					Login:    "my_login",
					Password: "my_password",
				})
				assert.Equal(t, ErrLoginExists, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}

func TestDBStorage_LoginUser(t *testing.T) {
	storage := newDBStorage("", "")
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	storage.DB = db

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Login user with good credentials",
			func() {
				mock.ExpectQuery(
					`SELECT user_id FROM users WHERE login = $1 AND password = $2`,
				).WithArgs("my_login", "my_password").WillReturnRows(
					sqlmock.NewRows([]string{"user_id"}).AddRow("12345678-1234-5678-9123-123456789012"))
			},
			func() {
				userID, err := storage.LoginUser(userdata.UserCredentials{
					Login:    "my_login",
					Password: "my_password",
				})
				assert.NoError(t, err)
				assert.Equal(t, userdata.UserID("12345678-1234-5678-9123-123456789012"), userID)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Login user with good credentials, but DB will return error",
			func() {
				mock.ExpectQuery(
					`SELECT user_id FROM users WHERE login = $1 AND password = $2`,
				).WithArgs("my_login", "my_password").WillReturnError(errors.New("some DB error"))
			},
			func() {
				userID, err := storage.LoginUser(userdata.UserCredentials{
					Login:    "my_login",
					Password: "my_password",
				})
				assert.Equal(t, ErrUnknown, err)
				assert.Equal(t, userdata.UserID(""), userID)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Login user with bad credentials",
			func() {
				mock.ExpectQuery(
					`SELECT user_id FROM users WHERE login = $1 AND password = $2`,
				).WithArgs("my_login", "my_password").WillReturnRows(sqlmock.NewRows([]string{"user_id"}))
			},
			func() {
				userID, err := storage.LoginUser(userdata.UserCredentials{
					Login:    "my_login",
					Password: "my_password",
				})
				assert.Equal(t, ErrWrongCredentials, err)
				assert.Equal(t, userdata.UserID(""), userID)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}

func TestDBStorage_GetRecordsInfo(t *testing.T) {
	storage := newDBStorage("", "")
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	storage.DB = db

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Get all info from unauthorized user",
			func() {},
			func() {
				records, err := storage.GetRecordsInfo(context.Background())
				assert.Equal(t, ErrUnauthenticated, err)
				assert.Empty(t, records)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Get all info from authorized user",
			func() {
				mock.ExpectQuery(
					"SELECT record_id, record_type, keyhint, metadata FROM data WHERE user_id = $1",
				).WithArgs("11111111-2222-33333-4444-555555555").WillReturnRows(
					sqlmock.NewRows([]string{"record_id", "record_type", "keyhint", "metadata"}).AddRow("1", userdata.TypeLoginAndPassword, "keyhint", "login and password").AddRow("2", userdata.TypeText, "keyhint", "custom text"))
			},
			func() {
				md := metadata.Pairs("userID", string("11111111-2222-33333-4444-555555555"))
				ctx := metadata.NewIncomingContext(context.Background(), md)

				records, err := storage.GetRecordsInfo(ctx)
				assert.NoError(t, err)

				assert.Equal(t, []userdata.Record{
					{
						ID:       "1",
						Type:     userdata.TypeLoginAndPassword,
						KeyHint:  "keyhint",
						Metadata: "login and password",
					},
					{
						ID:       "2",
						Type:     userdata.TypeText,
						KeyHint:  "keyhint",
						Metadata: "custom text",
					},
				}, records)

				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Get all info from authorized user, but DB will return error",
			func() {
				mock.ExpectQuery(
					"SELECT record_id, record_type, keyhint, metadata FROM data WHERE user_id = $1",
				).WithArgs(
					"11111111-2222-33333-4444-555555555",
				).WillReturnError(errors.New("some DB error"))
			},
			func() {
				md := metadata.Pairs("userID", string("11111111-2222-33333-4444-555555555"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
		
				records, err := storage.GetRecordsInfo(ctx)
				assert.Equal(t, ErrUnknown, err)
				assert.Empty(t, records)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}

func TestDBStorage_CreateRecord(t *testing.T) {
	storage := newDBStorage("", "")
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	storage.DB = db

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Create record with unauthorized user",
			func() {},
			func() {
				recordID, err := storage.CreateRecord(context.Background(), userdata.Record{})
				assert.Equal(t, ErrUnauthenticated, err)
				assert.Empty(t, recordID)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Create record with authorized user",
			func() {
				mock.ExpectQuery(
					"INSERT INTO data (user_id, record_type, keyhint, metadata, crypted_data) VALUES ($1, $2, $3, $4, $5) RETURNING record_id",
				).WithArgs(
					"11111111-2222-33333-4444-555555555",
					userdata.TypeText,
					"keyhint",
					"my text",
					hex.EncodeToString([]byte("hello!")),
				).WillReturnRows(sqlmock.NewRows([]string{"record_id"}).AddRow("1"))
			},
			func() {
				md := metadata.Pairs("userID", string("11111111-2222-33333-4444-555555555"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
				
				recordID, err := storage.CreateRecord(ctx, userdata.Record{
					KeyHint:  "keyhint",
					Metadata: "my text",
					Type:     userdata.TypeText,
					Data:     []byte("hello!"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", recordID)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Create record with authorized user, but DB will return error",
			func() {
				mock.ExpectQuery(
					"INSERT INTO data (user_id, record_type, keyhint, metadata, crypted_data) VALUES ($1, $2, $3, $4, $5) RETURNING record_id",
				).WithArgs(
					"11111111-2222-33333-4444-555555555",
					userdata.TypeText,
					"keyhint",
					"my text",
					hex.EncodeToString([]byte("hello!")),
				).WillReturnError(errors.New("some DB error"))
			},
			func() {
				md := metadata.Pairs("userID", string("11111111-2222-33333-4444-555555555"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
				
				recordID, err := storage.CreateRecord(ctx, userdata.Record{
					KeyHint:  "keyhint",
					Metadata: "my text",
					Type:     userdata.TypeText,
					Data:     []byte("hello!"),
				})
				assert.Equal(t, ErrUnknown, err)
				assert.Empty(t, recordID)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}

func TestDBStorage_GetRecord(t *testing.T) {
	storage := newDBStorage("", "")
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	storage.DB = db

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Get record with unauthorized user",
			func() {},
			func() {
				record, err := storage.GetRecord(context.Background(), "1")
				assert.Equal(t, ErrUnauthenticated, err)
				assert.Empty(t, record)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Get record with authorized user",
			func() {
				mock.ExpectQuery(
					"SELECT record_id, record_type, keyhint, metadata, crypted_data FROM data WHERE record_id = $1 AND user_id = $2",
				).WithArgs(
					"1", "11111111-2222-33333-4444-555555555",
				).WillReturnRows(
					sqlmock.NewRows([]string{"record_id", "record_type", "keyhint", "metadata", "crypted_data"}).AddRow("1", userdata.TypeText, "keyhint", "my text", hex.EncodeToString([]byte("hello!"))))
			},
			func() {
				md := metadata.Pairs("userID", string("11111111-2222-33333-4444-555555555"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
				
				record, err := storage.GetRecord(ctx, "1")
				assert.NoError(t, err)
				assert.Equal(t, userdata.Record{
					ID:       "1",
					Metadata: "my text",
					KeyHint:  "keyhint",
					Type:     userdata.TypeText,
					Data:     []byte("hello!"),
				}, record)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Get non existed record with authorized user",
			func() {
				mock.ExpectQuery(
					"SELECT record_id, record_type, keyhint, metadata, crypted_data FROM data WHERE record_id = $1 AND user_id = $2",
				).WithArgs(
					"1", "11111111-2222-33333-4444-555555555",
				).WillReturnRows(sqlmock.NewRows([]string{"record_id", "record_type", "keyhint", "metadata", "crypted_data"}))
			},
			func() {
				md := metadata.Pairs("userID", string("11111111-2222-33333-4444-555555555"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
				
				record, err := storage.GetRecord(ctx, "1")
				assert.Equal(t, ErrNotFound, err)
				assert.Empty(t, record)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Get record with authorized user, but DB will return error",
			func() {
				mock.ExpectQuery(
					"SELECT record_id, record_type, keyhint, metadata, crypted_data FROM data WHERE record_id = $1 AND user_id = $2",
				).WithArgs(
					"1", "11111111-2222-33333-4444-555555555",
				).WillReturnError(errors.New("some DB error"))
			},
			func() {
				md := metadata.Pairs("userID", string("11111111-2222-33333-4444-555555555"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
				
				record, err := storage.GetRecord(ctx, "1")
				assert.Equal(t, ErrUnknown, err)
				assert.Empty(t, record)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}

func TestDBStorage_DeleteRecord(t *testing.T) {
	storage := newDBStorage("", "")
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	storage.DB = db

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Delete record with unauthorized user",
			func() {},
			func() {
				err := storage.DeleteRecord(context.Background(), "1")
				assert.Equal(t, ErrUnauthenticated, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Delete record with authorized user",
			func() {
				mock.ExpectExec(
					"DELETE FROM data WHERE record_id = $1 AND user_id = $2",
				).WithArgs(
					"1", "11111111-2222-33333-4444-555555555",
				).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			func() {
				md := metadata.Pairs("userID", string("11111111-2222-33333-4444-555555555"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
				
				err := storage.DeleteRecord(ctx, "1")
				assert.NoError(t, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Delete record with authorized user, but DB will return error",
			func() {
				mock.ExpectExec(
					"DELETE FROM data WHERE record_id = $1 AND user_id = $2",
				).WithArgs(
					"1", "11111111-2222-33333-4444-555555555",
				).WillReturnError(errors.New("some DB error"))
			},
			func() {
				md := metadata.Pairs("userID", string("11111111-2222-33333-4444-555555555"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
			
				err := storage.DeleteRecord(ctx, "1")
				assert.Equal(t, ErrUnknown, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			"Delete non existed record with authorized user",
			func() {
				mock.ExpectExec(
					"DELETE FROM data WHERE record_id = $1 AND user_id = $2",
				).WithArgs(
					"1", "11111111-2222-33333-4444-555555555",
				).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			func() {
				md := metadata.Pairs("userID", string("11111111-2222-33333-4444-555555555"))
				ctx := metadata.NewIncomingContext(context.Background(), md)
				
				err := storage.DeleteRecord(ctx, "1")
				assert.Equal(t, ErrNotFound, err)
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
	}
}
