package clientconfig

import (
	"flag"
	"os"
)

// ClientConfig struct for client config.
type ClientConfig struct {
	ServerAddress string
	ClientCert    string
}

var (
	defaultServerAddress = "127.0.0.1:9000"
	defaultClientCert    = "../../cmd/cert/ca-cert.pem"
)

func NewClientConfig() ClientConfig {
	var cfg ClientConfig

	flag.StringVar(&cfg.ServerAddress, "addr", defaultServerAddress, "Server address and port")
	flag.StringVar(&cfg.ClientCert, "clientcert", defaultClientCert, "Path to client certificat for TLS")

	flag.Parse()

	if v, ok := os.LookupEnv("SERVER_ADDR"); ok {
		cfg.ServerAddress = v
	}

	if v, ok := os.LookupEnv("CLIENT_CERT"); ok {
		cfg.ClientCert = v
	}

	return cfg
}
