package estimator

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type Estimator interface {
	EstimateLastChlorine(ctx context.Context, poolID uuid.UUID) (float64, error)
	EstimateChlorine(ctx context.Context, poolID uuid.UUID, calciumHypochlorite65Percent, sodiumHypochlorite12Percent, sodiumHypochlorite14Percent, tCCA90PercentTablets, multiActionTablets, tCCA90PercentGranules, dichlor65Percent null.Float) (float64, error)
}

type repo interface {
	GetPool(ctx context.Context, id uuid.UUID) (*common.Pool, error)
}

type historyRepo interface {
	QueryMeasurement(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Measurement, error)
	QueryAdditives(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Additives, error)
}

type estimator struct {
	repo
	historyRepo

	log log.Logger
}

func (e *estimator) EstimateChlorine(ctx context.Context, poolID uuid.UUID, calciumHypochlorite65Percent, sodiumHypochlorite12Percent, sodiumHypochlorite14Percent, tCCA90PercentTablets, multiActionTablets, tCCA90PercentGranules, dichlor65Percent null.Float) (float64, error) {
	pool, err := e.GetPool(ctx, poolID)
	if err != nil {
		return 0, err
	}

	lastMeasurement, err := e.historyRepo.QueryMeasurement(ctx, poolID, common.OrderDesc)
	if err != nil {
		return 0, err
	}

	if len(lastMeasurement) == 0 {
		return 0, nil
	}

	additives := map[common.ChemicalProduct]float64{}

	if calciumHypochlorite65Percent.Valid {
		additives[common.CalciumHypochlorite65Percent] = calciumHypochlorite65Percent.Float64
	}

	if sodiumHypochlorite12Percent.Valid {
		additives[common.SodiumHypochlorite12Percent] = sodiumHypochlorite12Percent.Float64
	}

	if sodiumHypochlorite14Percent.Valid {
		additives[common.SodiumHypochlorite14Percent] = sodiumHypochlorite14Percent.Float64
	}

	if tCCA90PercentTablets.Valid {
		additives[common.TCCA90PercentTablets] = tCCA90PercentTablets.Float64
	}

	if multiActionTablets.Valid {
		additives[common.MultiActionTablets] = multiActionTablets.Float64
	}

	if tCCA90PercentGranules.Valid {
		additives[common.TCCA90PercentGranules] = tCCA90PercentGranules.Float64
	}

	if dichlor65Percent.Valid {
		additives[common.Dichlor65Percent] = dichlor65Percent.Float64
	}

	return CalculateChlorine(pool.Volume, lastMeasurement[0], additives)
}

func (e *estimator) EstimateLastChlorine(ctx context.Context, poolID uuid.UUID) (float64, error) {
	pool, err := e.GetPool(ctx, poolID)
	if err != nil {
		return 0, err
	}

	lastMeasurement, err := e.historyRepo.QueryMeasurement(ctx, poolID, common.OrderDesc)
	if err != nil {
		return 0, err
	}

	if len(lastMeasurement) == 0 {
		return 0, nil
	}

	lastAdditives, err := e.historyRepo.QueryAdditives(ctx, poolID, common.OrderDesc)
	if err != nil {
		return 0, err
	}

	if len(lastAdditives) == 0 {
		return lastMeasurement[0].Chlorine.Float64, nil
	}

	return CalculateChlorine(pool.Volume, lastMeasurement[0], lastAdditives[0].Products)
}

func NewEstimator(repo repo, historyRepo historyRepo, logger log.Logger) Estimator {
	return &estimator{repo: repo, historyRepo: historyRepo, log: logger}
}
