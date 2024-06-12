package userdata

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	TypeLoginAndPassword RecordType = iota
	TypeFile
	TypeText
	TypeCreditCard
)

// UserCredentials struct for user authorization.
type UserCredentials struct {
	Login    string
	Password string
	AESKey   string
}

// UserID is unique identificator.
type UserID string

// AuthToken is authorization token.
type AuthToken string

// Record is struct for decrypted or encrypted information.
type Record struct {
	ID       string
	Metadata string
	KeyHint  string
	Type     RecordType
	Data     []byte
}

type RecordType int32

func (r RecordType) String() string {
	switch r {
	case TypeLoginAndPassword:
		return "Login & password"
	case TypeFile:
		return "File"
	case TypeCreditCard:
		return "Credit card"
	case TypeText:
		return "Text"
	default:
		return "Unknown"
	}
}

// LoginAndPassword for encrypted login and password.
type LoginAndPassword struct {
	Login    string
	Password string
}

// Bytes implementation of Data interface.
func (lpdata *LoginAndPassword) Bytes() ([]byte, error) {
	return []byte(lpdata.Login + ":" + lpdata.Password), nil
}

// TextData for encrypted text data.
type TextData struct {
	Text string
}

// Bytes gets bytes.
func (tdata *TextData) Bytes() ([]byte, error) {
	return []byte(tdata.Text), nil
}

// BinaryFile for encrypted file.
type BinaryFile struct {
	FilePath string
	File     *os.File
}

// Bytes gets bytes.
func (fdata *BinaryFile) Bytes() ([]byte, error) {
	file, err := os.Open(fdata.FilePath)
	if err != nil {
		log.Infoln(err)

		return nil, err
	}
	fdata.File = file

	return io.ReadAll(fdata.File)
}

// CreditCard for credit card data.
type CreditCard struct {
	Number         string
	ExpirationDate string
	CVC            string
}

// Bytes gets bytes.
func (cdata *CreditCard) Bytes() ([]byte, error) {
	return []byte(cdata.Number + "/t" + cdata.ExpirationDate + "/t" + cdata.CVC), nil
}
