package additiveshistory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type AdditivesHistory interface {
	CreateChemicals(ctx context.Context, poolID uuid.UUID, r map[common.ChemicalProduct]float64) (*common.Chemicals, error)
	QueryAdditives(ctx context.Context, poolID uuid.UUID, order common.Order, offset, limit *int) ([]common.Chemicals, error)
	DeleteChemicals(ctx context.Context, poolID uuid.UUID, createdAt time.Time) (bool, error)
}

type repo interface {
	CreateChemicals(ctx context.Context, rec *common.Chemicals) error
	QueryChemicals(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Chemicals, error)
	QueryChemicalsGroupedByDay(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Chemicals, error)
}

type additivesHistory struct {
	repo

	log log.Logger
}

func (a *additivesHistory) DeleteChemicals(ctx context.Context, poolID uuid.UUID, createdAt time.Time) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (a *additivesHistory) CreateChemicals(
	ctx context.Context,
	poolID uuid.UUID,
	products map[common.ChemicalProduct]float64,
) (*common.Chemicals, error) {
	ch := &common.Chemicals{
		PoolID:    poolID,
		Products:  products,
		CreatedAt: time.Now(),
	}

	if err := a.repo.CreateChemicals(ctx, ch); err != nil {
		return nil, err
	}

	return ch, nil
}

func (a *additivesHistory) QueryAdditives(ctx context.Context, poolID uuid.UUID, order common.Order, offset, limit *int) ([]common.Chemicals, error) {
	return a.repo.QueryChemicals(ctx, poolID, order)
}

func NewAdditivesHistory(repo repo, log log.Logger) AdditivesHistory {
	return &additivesHistory{repo: repo, log: log}
}
