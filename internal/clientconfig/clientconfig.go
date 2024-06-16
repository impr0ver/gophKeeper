package clientconfig

import (
	"flag"
	"os"
	"strconv"
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
)

// ClientConfig struct for client config.
type ClientConfig struct {
	ServerAddress string
	ClientCert    string
	MaxFileSize   int64
}

var (
	defaultServerAddress = "127.0.0.1:9000"
	defaultClientCert    = "../../cmd/cert/ca-cert.pem"
	defaultMaxFileSize   = int64(8 * MB)
)

func NewClientConfig() ClientConfig {
	var cfg ClientConfig

	flag.StringVar(&cfg.ServerAddress, "addr", defaultServerAddress, "Server address and port")
	flag.StringVar(&cfg.ClientCert, "clientcert", defaultClientCert, "Path to client certificat for TLS")
	flag.Int64Var(&cfg.MaxFileSize, "maxsize", defaultMaxFileSize, "Max size of send file (type file) in MB")

	flag.Parse()

	if v, ok := os.LookupEnv("SERVER_ADDR"); ok {
		cfg.ServerAddress = v
	}

	if v, ok := os.LookupEnv("CLIENT_CERT"); ok {
		cfg.ClientCert = v
	}

	if v, ok := os.LookupEnv("FILE_MAXSIZE"); ok {
		int64Var, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			cfg.MaxFileSize = defaultMaxFileSize
		}
		cfg.MaxFileSize = int64Var * MB
	}

	return cfg
}
