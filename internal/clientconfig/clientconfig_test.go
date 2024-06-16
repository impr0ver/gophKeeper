package clientconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	os.Setenv("SERVER_ADDR", "127.0.0.1:9000")
	os.Setenv("FILE_MAXSIZE", "10")
	cfgTest := NewClientConfig()
	assert.Equal(t, "127.0.0.1:9000", cfgTest.ServerAddress, "test #SERVER_ADDR")
	os.Unsetenv("SERVER_ADDR")

	assert.Equal(t, int64(10*MB), cfgTest.MaxFileSize, "test #MaxFileSize")
	assert.Equal(t, "../../cmd/cert/ca-cert.pem", cfgTest.ClientCert, "test #ClientCert")

	assert.Equal(t, int64(10*MB), int64(10485760), "test #MaxFileSize2")
	os.Unsetenv("FILE_MAXSIZE")
}
