package additiveshistory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type AdditivesHistory interface {
	CreateAdditives(ctx context.Context, r *common.Additives) (*common.Additives, error)
	QueryAdditives(ctx context.Context, poolID uuid.UUID, order common.Order, offset, limit *int) ([]common.Additives, error)
}

type repo interface {
	CreateAdditive(ctx context.Context, rec *common.Additives) error
	QueryAdditives(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Additives, error)
}

type additivesHistory struct {
	repo

	log log.Logger
}

func (a *additivesHistory) CreateAdditives(ctx context.Context, r *common.Additives) (*common.Additives, error) {
	r.CreatedAt = time.Now()

	if err := a.repo.CreateAdditive(ctx, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (a *additivesHistory) QueryAdditives(ctx context.Context, poolID uuid.UUID, order common.Order, offset, limit *int) ([]common.Additives, error) {
	return a.repo.QueryAdditives(ctx, poolID, order)
}

func NewAdditivesHistory(repo repo, log log.Logger) AdditivesHistory {
	return &additivesHistory{repo: repo, log: log}
}
