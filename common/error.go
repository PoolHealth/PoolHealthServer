package common

import "errors"

var (
	ErrPoolNotFound        = errors.New("pool not found")
	ErrInvalidPoolID       = errors.New("invalid pool id")
	ErrPoolAlreadyExists   = errors.New("pool already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrExpiredToken        = errors.New("token already expired")
	ErrInvalidToken        = errors.New("token invalid")
	ErrJwtIncorrect        = errors.New("invalid jwt")
	ErrDeviceNotFound      = errors.New("device not found")
	ErrRecordAlreadyExists = errors.New("record already exists")
)
