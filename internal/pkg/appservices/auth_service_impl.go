package appservices

import (
	"crypto/rsa"
	"errors"
	"strconv"

	"github.com/Christian-007/fit-forge-notification/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge-notification/internal/pkg/cache"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"userId"`
	Uuid   string
	jwt.RegisteredClaims
}

type AuthData struct {
	UserId int `json:"userId"`
}

type AuthServiceImpl struct {
	AuthServiceOptions
}

type AuthServiceOptions struct {
	Cache cache.Cache
}

func NewAuthServiceImpl(options AuthServiceOptions) AuthServiceImpl {
	return AuthServiceImpl{
		options,
	}
}

func (a AuthServiceImpl) ValidateToken(privateKey *rsa.PrivateKey, tokenString string) (*Claims, error) {
	publicKey := &privateKey.PublicKey
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return publicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, apperrors.ErrInvalidSignature
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apperrors.ErrExpiredToken
		}

		return nil, err
	}

	if !token.Valid {
		return nil, apperrors.ErrInvalidToken
	}

	return claims, nil
}

func (a AuthServiceImpl) GetHashAuthDataFromCache(accessTokenUuid string) (AuthData, error) {
	result, err := a.Cache.GetAllHashFields(accessTokenUuid)
	if len(result) == 0 {
		return AuthData{}, apperrors.ErrRedisValueNotInHash
	}

	if err != nil {
		return AuthData{}, err
	}

	userIdInt, err := strconv.Atoi(result["userId"])
	if err != nil {
		return AuthData{}, err
	}

	return AuthData{
		UserId: userIdInt,
	}, nil
}
