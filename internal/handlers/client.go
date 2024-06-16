package handlers

import (
	"os"
	"sync"

	"github.com/impr0ver/gophKeeper/internal/crypt"
	"github.com/impr0ver/gophKeeper/internal/masker"
	"github.com/impr0ver/gophKeeper/internal/storage"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	log "github.com/sirupsen/logrus"
)

// client struct for client handlers.
type client struct {
	conn      ClientConnection
	authToken userdata.AuthToken
	AESKey    string
	Mu        *sync.Mutex
}

// newClientHandlers returns new client handlers with mutex.
func newClientHandlers(connection ClientConnection) *client {
	return &client{
		conn: connection,
		Mu:   &sync.Mutex{},
	}
}

// Login logins user by creds.
func (c *client) Login(credentials userdata.UserCredentials) error {
	if credentials.Login == "" || credentials.Password == "" || len(credentials.AESKey) == 0 {
		return ErrEmptyField
	}
	authToken, err := c.conn.Login(credentials)
	if err != nil {
		log.Warnf("%s :: %v", "auth token error", err)

		return err
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.authToken = userdata.AuthToken(authToken)
	c.AESKey = credentials.AESKey

	return nil
}

// SetAESKey reset the new AES key
func (c *client) SetAESKey(newAESKey string) error {
	if newAESKey == "" || len(newAESKey) == 0 {
		return ErrEmptyField
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.AESKey = newAESKey

	return nil
}

// Register creates new user by creds.
func (c *client) Register(credentials userdata.UserCredentials) error {
	if credentials.Login == "" || credentials.Password == "" || len(credentials.AESKey) == 0 {
		return ErrEmptyField
	}
	authToken, err := c.conn.Register(credentials)
	if err != nil {
		log.Warnf("%s :: %v", "register error", err)

		return err
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.authToken = userdata.AuthToken(authToken)
	c.AESKey = credentials.AESKey

	return nil
}

// GetRecordsInfo gets all records.
func (c *client) GetRecordsInfo() ([]userdata.Record, error) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	return c.conn.GetRecordsInfo(c.authToken)
}

// GetRecord gets record by recordID and decrypt cipherdata.
func (c *client) GetRecord(recordID string) (userdata.Record, error) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	record, errGetRecord := c.conn.GetRecord(c.authToken, recordID)
	if errGetRecord != nil {
		log.Infoln(errGetRecord)

		return record, errGetRecord
	}

	decrypted, err := crypt.AES256CBCDecode(record.Data, string(c.AESKey))
	if err != nil {
		log.Infoln(err)
		return record, storage.ErrUnknown
	}

	record.Data = decrypted

	// Get the file data and put in file
	if record.Type == userdata.TypeFile {
		file, err := os.Create(record.Metadata)
		if err != nil {
			log.Warnf("%s :: %v", "create file error", err)

			return record, storage.ErrUnknown
		}

		_, err = file.Write(record.Data)
		if err != nil {
			log.Warnf("%s :: %v", "write in file error", err)

			return record, storage.ErrUnknown
		}
		record.Data = []byte("Saved file successfully to " + record.Metadata + ".")
	}

	return record, nil
}

// DeleteRecord deletes record by ID.
func (c *client) DeleteRecord(recordID string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	return c.conn.DeleteRecord(c.authToken, recordID)
}

// CreateRecord creates new record and crypt plaindata.
func (c *client) CreateRecord(record userdata.Record) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	encrypted, err := crypt.AES256CBCEncode(record.Data, string(c.AESKey))
	if err != nil {
		log.Infoln(err)
		return storage.ErrUnknown
	}
	record.Data = encrypted

	//Add AES-key hint for facilities
	record.KeyHint = c.AESKey
	record.KeyHint = masker.Masker(record.KeyHint)

	return c.conn.CreateRecord(c.authToken, record)
}
