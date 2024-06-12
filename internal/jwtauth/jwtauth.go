package jwtauth

import (
	"github.com/impr0ver/gophKeeper/internal/userdata"
	"github.com/impr0ver/gophKeeper/internal/storage"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

// authenticatorJWT is authenticator which uses JWT.
type authenticatorJWT struct {
	secretKey      []byte
	expirationTime time.Duration
}

// NewAuthenticatorJWT gets new authenticatorJWT.
func NewAuthenticatorJWT(secretKey []byte, expirationTime time.Duration) *authenticatorJWT {
	return &authenticatorJWT{
		secretKey:      secretKey,
		expirationTime: expirationTime,
	}
}

// CreateToken implementation of Authenticator interface. Creates token, which stores userID.
func (a *authenticatorJWT) CreateToken(userID userdata.UserID) (userdata.AuthToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(a.expirationTime).Unix()
	claims["userID"] = userID

	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		log.Println("Failed generate token for authentication:", err)

		return "", storage.ErrUnknown
	}

	return userdata.AuthToken(tokenString), nil
}

// ValidateToken implementation of Authenticator interface. Validates token, returns userID.
func (a *authenticatorJWT) ValidateToken(token userdata.AuthToken) (userdata.UserID, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(string(token), claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, storage.ErrUnknown
		}
		return a.secretKey, nil
	})

	if err != nil {
		log.Warning(storage.ErrUnauthenticated)

		return "", storage.ErrUnauthenticated
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", storage.ErrUnauthenticated
	}

	return userdata.UserID(userID), nil
}
