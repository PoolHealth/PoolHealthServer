package estimator

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type repo interface {
	GetPool(ctx context.Context, id uuid.UUID) (*common.Pool, error)
}

type historyRepo interface {
	QueryMeasurement(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Measurement, error)
	QueryLastMeasurement(ctx context.Context, poolID uuid.UUID) ([]common.Measurement, error)
	QueryChemicalsGroupedByDay(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Chemicals, error)
}

type Estimator struct {
	repo
	historyRepo

	log log.Logger
}

func (e *Estimator) RecommendedChemicals(ctx context.Context, poolID uuid.UUID) (map[common.ChemicalProduct]float64, error) {
	//TODO implement me
	panic("implement me")
}

// TODO separate calculation
func (e *Estimator) DemandMeasurement(ctx context.Context, poolID uuid.UUID) (common.Measurement, error) {
	pool, err := e.GetPool(ctx, poolID)
	if err != nil {
		return common.Measurement{}, err
	}

	lastMeasurement, err := e.historyRepo.QueryMeasurement(ctx, poolID, common.OrderDesc)
	if err != nil {
		return common.Measurement{}, err
	}

	if len(lastMeasurement) == 0 {
		return common.Measurement{}, nil
	}

	lastAdditives, err := e.historyRepo.QueryChemicalsGroupedByDay(ctx, poolID, common.OrderDesc)
	if err != nil {
		return common.Measurement{}, err
	}

	chlorine := lastMeasurement[0].Chlorine.Float64
	if len(lastMeasurement) > 1 {
		chlorine -= lastMeasurement[1].Chlorine.Float64
	}
	ph := lastMeasurement[0].Chlorine.Float64
	alkalinity := lastMeasurement[0].Alkalinity.Float64

	if len(lastAdditives) != 0 {
		chlorine, err = CalculateChlorine(pool.Volume, lastMeasurement[0], lastAdditives[0].Products)
		if err != nil {
			return common.Measurement{}, err
		}

		if len(lastMeasurement) > 1 {
			chlorine -= lastMeasurement[1].Chlorine.Float64
		}
	}

	result := common.Measurement{
		PoolID:     poolID,
		Chlorine:   null.FloatFrom(chlorine),
		PH:         null.FloatFrom(ph),
		Alkalinity: null.FloatFrom(alkalinity),
	}

	return result, nil
}

func (e *Estimator) EstimateMeasurement(ctx context.Context, chemicals map[common.ChemicalProduct]float64) (common.Measurement, error) {
	//TODO implement me
	panic("implement me")
}

func (e *Estimator) EstimateChlorine(ctx context.Context, poolID uuid.UUID, calciumHypochlorite65Percent, sodiumHypochlorite12Percent, sodiumHypochlorite14Percent, tCCA90PercentTablets, multiActionTablets, tCCA90PercentGranules, dichlor65Percent null.Float) (float64, error) {
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

func NewEstimator(repo repo, historyRepo historyRepo, logger log.Logger) *Estimator {
	return &Estimator{repo: repo, historyRepo: historyRepo, log: logger}
}
