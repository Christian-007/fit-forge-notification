package apperrors

import "errors"

var (
	ErrInvalidSignature    = errors.New("invalid token signature")
	ErrInvalidToken        = errors.New("invalid token")
	ErrExpiredToken        = errors.New("expired token")
	ErrRedisKeyNotFound    = errors.New("key does not exist in Redis")
	ErrRedisValueNotInHash = errors.New("value is not in Redis Hash")
	ErrTypeAssertion       = errors.New("type assertion failed")
)
