package appservices

import (
	"crypto/rsa"
)

//go:generate mockgen -source=interface.go -destination=mocks/interface.go
type AuthService interface {
	ValidateToken(privateKey *rsa.PrivateKey, tokenString string) (*Claims, error)
	GetHashAuthDataFromCache(accessTokenUuid string) (AuthData, error)
}
