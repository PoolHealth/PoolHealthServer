package graphql

import (
	"context"

	"github.com/google/uuid"

	"github.com/PoolHealth/PoolHealthServer/common"
)

//go:generate go run ../scripts/gqlgen.go

type logger interface {
	Debug(msgs ...interface{})
}

type poolData interface {
	Create(ctx context.Context, userID uuid.UUID, data *common.PoolData) (pool *common.Pool, err error)
	Update(ctx context.Context, id uuid.UUID, rec *common.PoolData) (record *common.Pool, err error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID) ([]common.Pool, error)
	SubscribeOnCreate(ctx context.Context) (<-chan *common.Pool, error)
	SubscribeOnUpdate(ctx context.Context) (<-chan *common.Pool, error)
	SubscribeOnDelete(ctx context.Context) (<-chan uuid.UUID, error)
}

type auth interface {
	Auth(ctx context.Context, token string) (*common.Session, error)
}

type Resolver struct {
	poolData
	auth

	log logger
}

func NewResolver(
	logger logger,
	data poolData,
	auth auth,
) *Resolver {
	return &Resolver{
		auth:     auth,
		poolData: data,
		log:      logger,
	}
}
