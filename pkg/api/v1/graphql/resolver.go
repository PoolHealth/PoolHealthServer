package graphql

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/PoolHealth/PoolHealthServer/common"
	authPkg "github.com/PoolHealth/PoolHealthServer/internal/auth"
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
	Has(ctx context.Context, id uuid.UUID, userID uuid.UUID) (bool, error)
}

type measurementHistory interface {
	CreateMeasurement(ctx context.Context, r common.Measurement) (common.Measurement, error)
	QueryMeasurement(ctx context.Context, poolID uuid.UUID, order common.Order, offset, limit *int) ([]common.Measurement, error)
}

type additivesHistory interface {
	CreateAdditives(ctx context.Context, r *common.Additives) (*common.Additives, error)
	QueryAdditives(ctx context.Context, poolID uuid.UUID, order common.Order, offset, limit *int) ([]common.Additives, error)
}

type estimator interface {
	EstimateChlorine(ctx context.Context, poolID uuid.UUID, calciumHypochlorite65Percent, sodiumHypochlorite12Percent, sodiumHypochlorite14Percent, tCCA90PercentTablets, multiActionTablets, tCCA90PercentGranules, dichlor65Percent null.Float) (float64, error)
	EstimateLastChlorine(ctx context.Context, poolID uuid.UUID) (float64, error)
}

type auth interface {
	Auth(ctx context.Context, token string) (*common.Session, error)
}

type Resolver struct {
	poolData
	measurementHistory
	additivesHistory
	auth
	estimator estimator

	log logger
}

func (r *Resolver) checkAccessToPool(ctx context.Context, poolID uuid.UUID) error {
	user, err := authPkg.GetUser(ctx)
	if err != nil {
		return err
	}

	ok, err := r.poolData.Has(ctx, poolID, user.ID)
	if err != nil {
		return err
	}

	if !ok {
		return gqlerror.Errorf("pool not found")
	}

	return nil
}

func NewResolver(
	logger logger,
	data poolData,
	measurementHistory measurementHistory,
	additivesHistory additivesHistory,
	estimator estimator,
	auth auth,
) *Resolver {
	return &Resolver{
		auth:               auth,
		poolData:           data,
		measurementHistory: measurementHistory,
		additivesHistory:   additivesHistory,
		estimator:          estimator,
		log:                logger,
	}
}
