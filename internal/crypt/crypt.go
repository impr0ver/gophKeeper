package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/impr0ver/gophKeeper/internal/userdata"

	log "github.com/sirupsen/logrus"
	"github.com/zenazn/pkcs7pad"
)

// GenerateRandom generates random bytes for encrypting.
func GenerateRand(size int) ([]byte, error) {
	buffer := make([]byte, size)
	_, err := rand.Read(buffer)
	if err != nil {
		log.Infoln(err)

		return nil, err
	}

	return buffer, nil
}

// PasswordHash make encryption string.
func PasswordHash(credentials userdata.UserCredentials) string {
	sha := sha256.New()
	sha.Write([]byte(credentials.Login + credentials.Password))

	return hex.EncodeToString(sha.Sum(nil))
}

// getMD5Hash (cipher key must be 32 chars long because block size is 16 bytes!)
func getMD5Hash(text string) []byte {
	hash := md5.Sum([]byte(text))
	return hash[:]
}

// AES256CBCEncode encrypt plain text string into cipher text string
func AES256CBCEncode(plainText []byte, key string) ([]byte, error) {
	bKey := getMD5Hash(key)

	padded := pkcs7pad.Pad(plainText, aes.BlockSize)

	if len(padded)%aes.BlockSize != 0 {
		err := fmt.Errorf("plainText has the wrong block size")
		return nil, err
	}

	block, err := aes.NewCipher(bKey)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(padded))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], padded)

	return cipherText, nil
}

// AES256CBCDecode decrypt cipher text string into plain text string
func AES256CBCDecode(cipherText []byte, key string) ([]byte, error) {
	bKey := getMD5Hash(key)

	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		panic("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		panic("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, err = pkcs7pad.Unpad(cipherText)
	if err != nil {
		return nil, fmt.Errorf("error in Unpad function: %w", err)
	}
	return cipherText, nil
}
