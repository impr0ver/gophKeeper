package jwtauth

import (
	"testing"
	"time"

	"github.com/impr0ver/gophKeeper/internal/userdata"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthenticatorJWT(t *testing.T) {
	auth := NewAuthenticatorJWT([]byte("mySuperSecretKey"), time.Duration(1 * time.Hour))
	assert.NotEmpty(t, auth)
}

func TestAuthenticatorJWT(t *testing.T) {
	auth := NewAuthenticatorJWT([]byte("mySuperSecretKey"), time.Duration(1 * time.Hour))

	userID := userdata.UserID("ID7777")

	token, err := auth.CreateToken(userID)
	assert.NoError(t, err)

	id, errValidate := auth.ValidateToken(token)
	assert.NoError(t, errValidate)
	assert.Equal(t, userID, id)
}
