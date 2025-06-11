package estimator

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

var (
	ErrNoSettings  = errors.New("can't recommend chemicals without settings")
	ErrNoUsageType = errors.New("can't recommend chemicals without usage type")
)

type repo interface {
	GetPool(ctx context.Context, id uuid.UUID) (*common.Pool, error)
	GetSettings(ctx context.Context, poolID uuid.UUID) (*common.PoolSettings, error)
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

type LastMeasurement struct {
	Chlorine         float64
	PreviousChlorine float64
	PH               float64
	Alkalinity       float64
}

func (l LastMeasurement) Empty() bool {
	return l.Chlorine == 0 && l.PH == 0 && l.Alkalinity == 0 && l.PreviousChlorine == 0
}

func (l LastMeasurement) Filled() bool {
	return l.Chlorine != 0 && l.PH != 0 && l.Alkalinity != 0 && l.PreviousChlorine != 0
}

func (e *Estimator) RecommendedChemicals(ctx context.Context, poolID uuid.UUID) (map[common.ChemicalProduct]float64, error) {
	pool, err := e.GetPool(ctx, poolID)
	if err != nil {
		return nil, err
	}

	settings, err := e.GetSettings(ctx, poolID)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		return nil, ErrNoSettings
	}

	if settings.UsageType == common.UsageTypeUnknown {
		return nil, ErrNoUsageType
	}

	highTarget := chlorineHighByUsage[settings.UsageType]

	lastMeasurements, err := e.getLastMeasurements(ctx, poolID)
	if err != nil {
		return nil, err
	}

	if lastMeasurements.Empty() {
		return nil, nil
	}

	data := recommendChlorineByTarget(pool.Volume, lastMeasurements.Chlorine, highTarget)

	for k, v := range data {
		data[k] = math.Trunc(v*100) / 100
	}

	return data, nil
}

// TODO separate calculation
func (e *Estimator) DemandMeasurement(ctx context.Context, poolID uuid.UUID) (common.Measurement, error) {
	pool, err := e.GetPool(ctx, poolID)
	if err != nil {
		return common.Measurement{}, err
	}

	lastMeasurements, err := e.getLastMeasurements(ctx, poolID)
	if err != nil {
		return common.Measurement{}, err
	}

	if lastMeasurements.Empty() {
		return common.Measurement{}, nil
	}

	lastAdditives, err := e.historyRepo.QueryChemicalsGroupedByDay(ctx, poolID, common.OrderDesc)
	if err != nil {
		return common.Measurement{}, err
	}

	chlorine := lastMeasurements.Chlorine
	ph := lastMeasurements.PH
	alkalinity := lastMeasurements.Alkalinity

	if len(lastAdditives) != 0 {
		chlorine = CalculateChlorine(pool.Volume, lastMeasurements.Chlorine, lastAdditives[0].Products)

		ph = CalculatePH(pool.Volume, lastMeasurements.PH, lastAdditives[0].Products)

		alkalinity = CalculateAlkalinity(pool.Volume, lastMeasurements.Alkalinity, lastAdditives[0].Products)
	}

	chlorine -= lastMeasurements.PreviousChlorine

	result := common.Measurement{
		PoolID:     poolID,
		Chlorine:   null.FloatFrom(chlorine),
		PH:         null.FloatFrom(ph),
		Alkalinity: null.FloatFrom(alkalinity),
	}

	return result, nil
}

func (e *Estimator) EstimateMeasurement(
	ctx context.Context,
	poolID uuid.UUID,
	chemicals map[common.ChemicalProduct]float64,
	selector []common.MeasurementType,
) (common.Measurement, error) {
	pool, err := e.GetPool(ctx, poolID)
	if err != nil {
		return common.Measurement{}, err
	}

	lastMeasurement, err := e.getLastMeasurements(ctx, poolID)
	if err != nil {
		return common.Measurement{}, err
	}

	if lastMeasurement.Empty() {
		return common.Measurement{}, nil
	}

	result := common.Measurement{}

	for _, s := range selector {
		switch s {
		case common.MeasurementChlorine:
			result.Chlorine = null.FloatFrom(CalculateChlorine(pool.Volume, lastMeasurement.Chlorine, chemicals))
		case common.MeasurementPH:
			result.PH = null.FloatFrom(CalculatePH(pool.Volume, lastMeasurement.PH, chemicals))
		case common.MeasurementAlkalinity:
			result.Alkalinity = null.FloatFrom(CalculateAlkalinity(pool.Volume, lastMeasurement.Alkalinity, chemicals))
		}
	}

	return result, nil
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

	return CalculateChlorine(pool.Volume, lastMeasurement[0].Chlorine.Float64, additives), nil
}

func (e *Estimator) getLastMeasurements(ctx context.Context, poolID uuid.UUID) (LastMeasurement, error) {
	lastMeasurement, err := e.historyRepo.QueryMeasurement(ctx, poolID, common.OrderDesc)
	if err != nil {
		return LastMeasurement{}, err
	}

	if len(lastMeasurement) == 0 {
		return LastMeasurement{}, nil
	}

	result := LastMeasurement{}

	for _, m := range lastMeasurement {
		if m.Chlorine.Valid {
			if result.Chlorine != 0 {
				result.PreviousChlorine = m.Chlorine.Float64
			} else {
				result.Chlorine = m.Chlorine.Float64
			}

		}

		if m.PH.Valid {
			result.PH = m.PH.Float64
		}

		if m.Alkalinity.Valid {
			result.Alkalinity = m.Alkalinity.Float64
		}

		if result.Filled() {
			return result, nil
		}
	}

	return result, nil
}

func NewEstimator(repo repo, historyRepo historyRepo, logger log.Logger) *Estimator {
	return &Estimator{repo: repo, historyRepo: historyRepo, log: logger}
}
