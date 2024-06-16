package storage

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/impr0ver/gophKeeper/internal/logger"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

// dbStorage for db storage.
type dbStorage struct {
	DB     *sql.DB
	MigURL string
}

// NewDBStorage connects to DB.
func newDBStorage(connectionURL string, migrateURL string) *dbStorage {
	db, err := sql.Open("pgx", connectionURL)
	if err != nil {
		log.Fatalln("Failed open DB storage:", err)

		return nil
	}

	return &dbStorage{
		DB:     db,
		MigURL: migrateURL,
	}
}

// MigrateUP migrates DB.
func (ds *dbStorage) MigrateUP() {
	var sLogger = logger.NewSugarLogger()
	driver, err := postgres.WithInstance(ds.DB, &postgres.Config{})
	if err != nil {
		log.Infof("Failed create postgres instance: %v\n", err)
		sLogger.Fatalf("Failed create postgres instance: %v\n", err)
	}

	mig, errMigrate := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", ds.MigURL), "pgx", driver)
	if errMigrate != nil {
		log.Infof("Failed create migration instance: %v\n", err)
		sLogger.Fatalf("Failed create migration instance: %v\n", err)
		return
	}

	// For migrate down if need
	// if err := mig.Down(); err != nil && err != migrate.ErrNoChange {
	// 	log.Fatalln("Failed migrate: ", err)
	// 	return
	// }

	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalln("Failed migrate: ", err)

		return
	}
}

// CreateUser saves to DB new user.
func (ds *dbStorage) CreateUser(credentials userdata.UserCredentials) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := ds.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE login = $1`, credentials.Login)

	sameLoginCounter := 0
	err := row.Scan(&sameLoginCounter)

	if err != nil || row.Err() != nil {
		log.Infoln(err)

		return ErrUnknown
	}

	if sameLoginCounter > 0 {
		return ErrLoginExists
	}

	_, err = ds.DB.ExecContext(ctx, `INSERT INTO users (login, password) VALUES ($1, $2)`, credentials.Login, credentials.Password)
	if err != nil {
		log.Infoln(err)

		return ErrUnknown
	}

	return nil
}

// LoginUser check if credentials are valid and return userID.
func (ds *dbStorage) LoginUser(credentials userdata.UserCredentials) (userdata.UserID, error) {
	var userID userdata.UserID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := ds.DB.QueryRowContext(ctx, `SELECT user_id FROM users WHERE login = $1 AND password = $2`, credentials.Login, credentials.Password)

	err := row.Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Infoln(err)

		return userID, ErrWrongCredentials
	}

	if err != nil || row.Err() != nil {
		log.Infoln(err)

		return userID, ErrUnknown
	}

	return userID, nil
}

// GetRecordsInfo gets all DB record by userID.
func (ds *dbStorage) GetRecordsInfo(ctx context.Context) ([]userdata.Record, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("userID")) == 0 {
		log.Println("Failed get userID from context in getting all records")
		return nil, ErrUnauthenticated
	}

	userID := userdata.UserID(md.Get("userID")[0])

	rows, err := ds.DB.QueryContext(ctx, `SELECT record_id, record_type, keyhint, metadata FROM data WHERE user_id = $1`, userID)
	if err != nil {
		log.Infoln(err)

		return nil, ErrUnknown
	}

	defer rows.Close()

	result := make([]userdata.Record, 0, 10)

	var row userdata.Record
	for rows.Next() {
		if err := rows.Scan(&row.ID, &row.Type, &row.KeyHint, &row.Metadata); err != nil {
			log.Infoln(err)

			return nil, ErrUnknown
		}

		result = append(result, row)
	}

	if rows.Err() != nil {
		log.Println("Failed get rows in getting all records:", err)
		return nil, ErrUnknown
	}

	return result, nil
}

// CreateRecord saves new record to DB and return recordID.
func (ds *dbStorage) CreateRecord(ctx context.Context, record userdata.Record) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("userID")) == 0 {
		log.Println("Failed get userID from context in getting all records")
		return "", ErrUnauthenticated
	}

	userID := userdata.UserID(md.Get("userID")[0])

	hexDataString := hex.EncodeToString(record.Data)

	row := ds.DB.QueryRowContext(ctx, `INSERT INTO data (user_id, record_type, keyhint, metadata, crypted_data) VALUES ($1, $2, $3, $4, $5) RETURNING record_id`,
		userID,
		record.Type,
		record.KeyHint,
		record.Metadata,
		hexDataString,
	)

	var recordID string
	if err := row.Scan(&recordID); err != nil || row.Err() != nil {
		log.Infoln(err)

		return "", ErrUnknown
	}

	return recordID, nil
}

// GetRecord gets record from DB by userID.
func (ds *dbStorage) GetRecord(ctx context.Context, recordID string) (userdata.Record, error) {
	record := userdata.Record{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("userID")) == 0 {
		log.Println("Failed get userID from context in getting all records")
		return record, ErrUnauthenticated
	}

	userID := userdata.UserID(md.Get("userID")[0])

	row := ds.DB.QueryRowContext(ctx, `SELECT record_id, record_type, keyhint, metadata, crypted_data FROM data WHERE record_id = $1 AND user_id = $2`,
		recordID,
		userID,
	)

	var hexDataString string
	err := row.Scan(&record.ID, &record.Type, &record.KeyHint, &record.Metadata, &hexDataString)

	if errors.Is(err, sql.ErrNoRows) {
		log.Infoln(err)

		return record, ErrNotFound
	}
	if err != nil || row.Err() != nil {
		log.Infoln(err)

		return record, ErrUnknown
	}

	record.Data, err = hex.DecodeString(hexDataString)
	if err != nil {
		log.Infoln(err)

		return record, ErrUnknown
	}

	return record, nil
}

// DeleteRecord deletes record from DB by userID.
func (ds *dbStorage) DeleteRecord(ctx context.Context, recordID string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("userID")) == 0 {
		log.Println("Failed get userID from context in getting all records")
		return ErrUnauthenticated
	}

	userID := userdata.UserID(md.Get("userID")[0])

	result, err := ds.DB.ExecContext(ctx, `DELETE FROM data WHERE record_id = $1 AND user_id = $2`, recordID, userID)
	if err != nil {
		log.Infoln(err)

		return ErrUnknown
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		log.Println("Failed get affected records:", err)
		return ErrUnknown
	} else if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
