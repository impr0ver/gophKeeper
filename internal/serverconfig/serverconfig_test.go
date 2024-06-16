package serverconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	os.Setenv("SERVER_PORT", "127.0.0.1:9000")
	os.Setenv("SERVCONS_LOG", "TRUE")
	os.Setenv("DATABASE_DSN", "user=postgres password=ggggg host=localhost port=5432 dbname=gokeeper sslmode=disable")
	os.Setenv("FILE_STORAGE_PATH", "data")
	os.Setenv("SECRET_JWT", "mySuperSecretKey")
	os.Setenv("SERVER_CERT", "../../cmd/cert/server-cert.pem")
	os.Setenv("SERVER_KEY", "../../cmd/cert/server-key.pem")
	os.Setenv("EXP_TIME", "1s")
	os.Setenv("MIGRATE_URL", "../../migrations")

	cfgTest := NewServerConfig()

	assert.NotEmpty(t, cfgTest)
	assert.Equal(t, "127.0.0.1:9000", cfgTest.ListenAddr, "test #SERVER_ADDR")
	os.Unsetenv("SERVER_PORT")

	assert.Equal(t, "user=postgres password=ggggg host=localhost port=5432 dbname=gokeeper sslmode=disable", cfgTest.DatabaseDSN, "test #DatabaseDSN")
	assert.Equal(t, "data", cfgTest.FilesStore, "test #FilesStore")
	assert.Equal(t, "../../cmd/cert/server-key.pem", cfgTest.ServerKey, "test #ServerKey")

	assert.Equal(t, true, cfgTest.ServerConsoleLog, "test #ServerConsoleLog")
	os.Unsetenv("SERVCONS_LOG")
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("SECRET_JWT")
	os.Unsetenv("SERVER_CERT")
	os.Unsetenv("SERVER_KEY")
	os.Unsetenv("EXP_TIME")
	os.Unsetenv("MIGRATE_URL")
}
