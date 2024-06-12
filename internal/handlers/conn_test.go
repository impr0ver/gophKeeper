package handlers

import (
	"context"
	"testing"

	"github.com/impr0ver/gophKeeper/internal/clientconfig"
	"github.com/impr0ver/gophKeeper/internal/handlers/mocks"
	"github.com/impr0ver/gophKeeper/internal/serverconfig"
	"github.com/impr0ver/gophKeeper/internal/storage"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	var serverCfg = serverconfig.ServerConfig{}
	serverCfg.ServerCert = "../../cmd/cert/server-cert.pem"
	serverCfg.ServerKey = "../../cmd/cert/server-key.pem"
	serverCfg.ListenAddr = "127.0.0.1:9000"

	var clientCfg = clientconfig.ClientConfig{}
	clientCfg.ServerAddress = "127.0.0.1:9000"
	clientCfg.ClientCert = "../../cmd/cert/ca-cert.pem"

	auth := mocks.NewAuthenticator(t)
	handlers := mocks.NewServerHandlers(t)

	server := NewServerConn(handlers, auth, serverCfg.ServerCert, serverCfg.ServerKey, serverCfg.ServerConsoleLog)
	server.Run(context.Background(), serverCfg.ListenAddr)
	defer server.Stop()

	client := newClientConn(clientCfg.ServerAddress, clientCfg.ClientCert)

	tc := []struct {
		name  string
		mock  func()
		valid func()
	}{
		{
			"Create user",
			func() {
				handlers.On("CreateUser", userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				}).Return(userdata.AuthToken("token"), nil).Once()
			},
			func() {
				token, err := client.Register(userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				})
				assert.NoError(t, err)
				assert.Equal(t, "token", token)
			},
		},
		{
			"Create user, but error",
			func() {
				handlers.On("CreateUser", userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				}).Return(userdata.AuthToken("token"), storage.ErrLoginExists).Once()
			},
			func() {
				token, err := client.Register(userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				})
				assert.Equal(t, storage.ErrLoginExists, err)
				assert.Empty(t, token)
			},
		},
		{
			"Create user, but unknown error",
			func() {
				handlers.On("CreateUser", userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				}).Return(userdata.AuthToken("token"), storage.ErrUnknown).Once()
			},
			func() {
				token, err := client.Register(userdata.UserCredentials{
					Login:    "Login",
					Password: "Password",
				})
				assert.Equal(t, storage.ErrUnknown, err)
				assert.Empty(t, token)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.mock()
		test.valid()
		handlers.AssertExpectations(t)
	}
}
