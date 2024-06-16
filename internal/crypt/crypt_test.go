package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES(t *testing.T) {

	t.Run("Encrypts and decrypts", func(t *testing.T) {

		plainTexts := []string{"1234567890", "123456789012345678901234567890123456789012345678901234567890", "1", "MySecretText", ""}
		key := "mytestkey"
		for _, plainText := range plainTexts {

			encrypted, err := AES256CBCEncode([]byte(plainText), key)
			if err != nil {
				t.Fatalf("Failed to encrypt: %s - %s", []byte(plainText), err.Error())
			}

			decrypted, err := AES256CBCDecode(encrypted, key)
			if err != nil {
				t.Fatalf("Failed to decrypt: %s - %s", []byte(plainText), err.Error())
			}

			assert.Equal(t, []byte(plainText), decrypted)
			assert.Equal(t, len([]byte(plainText)), len(decrypted))
		}
	})
}

func TestAES_second(t *testing.T) {

	t.Run("Encrypts and decrypts", func(t *testing.T) {

		plainTexts := []string{"1234567890", "123456789012345678901234567890123456789012345678901234567890", "1", "MySecretText", ""}
		key := "mytestkey"
		key2 := "testkeymy"
		for _, plainText := range plainTexts {

			encrypted, err := AES256CBCEncode([]byte(plainText), key)
			assert.NoError(t, err)
			
			_, err = AES256CBCDecode(encrypted, key2)
			if err != nil {
				assert.Equal(t, "error in Unpad function: pkcs7pad: bad padding", err.Error())
			}
		}
	})
}

func Test_GenerateRandom(t *testing.T) {
	bytes, err := GenerateRand(32)
	assert.NoError(t, err)
	assert.NotEmpty(t, bytes)
	assert.Len(t, bytes, 32)
}
