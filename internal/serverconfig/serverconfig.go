package serverconfig

import (
	"flag"
	"os"
	"strconv"
	"time"
)

// ServerConfig struct for server config.
type ServerConfig struct {
	ListenAddr       string
	DatabaseDSN      string
	FilesStore       string
	JWTAuth          AuthConfig
	ServerCert       string
	ServerKey        string
	MigrationsURL    string
	ServerConsoleLog bool
}

// AuthConfig auth settings.
type AuthConfig struct {
	SecretJWT      string
	ExpirationTime time.Duration
}

var (
	defaultListenAddr       = "127.0.0.1:9000"
	defaultFilesStore       = "data"
	defaultDSN              = "user=postgres password=karat911 host=localhost port=5432 dbname=gokeeper sslmode=disable" //user=postgres password=mypassword host=localhost port=5432 dbname=gokeeper sslmode=disable
	defaultJWTSecret        = "mySuperSecretKey"
	defaultExpirationTime   = time.Duration(2 * time.Minute)
	defaultServerCert       = "../../cmd/cert/server-cert.pem"
	defaultServerKey        = "../../cmd/cert/server-key.pem"
	defaultMigrationsURL    = "../../migrations"
	defaultServerConsoleLog = true
)

// NewServerConfig gets server config.
func NewServerConfig() ServerConfig {
	var (
		cfg ServerConfig
		err error
	)

	flag.StringVar(&cfg.ListenAddr, "listenaddr", defaultListenAddr, "Server address and port")
	flag.StringVar(&cfg.DatabaseDSN, "dsn", defaultDSN, "Source to DB")
	flag.StringVar(&cfg.FilesStore, "filepath", defaultFilesStore, "Path to files store")
	flag.StringVar(&cfg.JWTAuth.SecretJWT, "jwtsecr", defaultJWTSecret, "JWT secret")
	flag.DurationVar(&cfg.JWTAuth.ExpirationTime, "exptime", defaultExpirationTime, "Token expiration time")
	flag.StringVar(&cfg.ServerCert, "servcert", defaultServerCert, "Path to server certificat for TLS")
	flag.StringVar(&cfg.ServerKey, "servkey", defaultServerKey, "Path to server key for TLS")
	flag.StringVar(&cfg.MigrationsURL, "migrateURL", defaultMigrationsURL, "Path to migrations for DB")
	flag.BoolVar(&cfg.ServerConsoleLog, "servconslog", defaultServerConsoleLog, "Console log request and MD data on server interceptors")

	flag.Parse()

	if v, ok := os.LookupEnv("SERVER_PORT"); ok {
		cfg.ListenAddr = v
	}

	if v, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.DatabaseDSN = v
	}

	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FilesStore = v
	}

	if v, ok := os.LookupEnv("SECRET_JWT"); ok {
		cfg.JWTAuth.SecretJWT = v
	}

	if v, ok := os.LookupEnv("EXP_TIME"); ok {
		cfg.JWTAuth.ExpirationTime, err = time.ParseDuration(v)
		if err != nil {
			cfg.JWTAuth.ExpirationTime = defaultExpirationTime
		}
	}

	if v, ok := os.LookupEnv("SERVER_CERT"); ok {
		cfg.ServerCert = v
	}

	if v, ok := os.LookupEnv("SERVER_KEY"); ok {
		cfg.ServerKey = v
	}

	if v, ok := os.LookupEnv("MIGRATE_URL"); ok {
		cfg.MigrationsURL = v
	}

	if v, ok := os.LookupEnv("SERVCONS_LOG"); ok {
		cfg.ServerConsoleLog, err = strconv.ParseBool(v)
		if err != nil {
			cfg.ServerConsoleLog = defaultServerConsoleLog
		}
	}

	return cfg
}
