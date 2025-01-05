package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/PoolHealth/PoolHealthServer/common"
)

type GQLErrorCode int

const (
	ErrExpiredToken GQLErrorCode = iota
	ErrInvalidToken
	ErrUserNotFound
	ErrPoolNotFound
	ErrDeviceNotFound
)

var errorsMap = map[error]GQLErrorCode{
	common.ErrExpiredToken:   ErrExpiredToken,
	common.ErrInvalidToken:   ErrInvalidToken,
	common.ErrUserNotFound:   ErrUserNotFound,
	common.ErrPoolNotFound:   ErrPoolNotFound,
	common.ErrDeviceNotFound: ErrDeviceNotFound,
}

func castGQLError(ctx context.Context, err error) error {
	extensions := map[string]interface{}{}

	code, ok := errorsMap[err]
	if ok {
		extensions["code"] = code
	}

	return &gqlerror.Error{
		Message:    err.Error(),
		Path:       graphql.GetPath(ctx),
		Extensions: extensions,
	}
}
