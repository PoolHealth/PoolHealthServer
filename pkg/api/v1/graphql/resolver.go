package graphql

import (
	"context"
	"time"

	gql "github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/PoolHealth/PoolHealthServer/common"
	authPkg "github.com/PoolHealth/PoolHealthServer/internal/services/auth"
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
	DeleteMeasurement(ctx context.Context, poolID uuid.UUID, createdAt time.Time) (bool, error)
}

type additivesHistory interface {
	CreateChemicals(ctx context.Context, poolID uuid.UUID, r map[common.ChemicalProduct]float64) (*common.Chemicals, error)
	QueryAdditives(ctx context.Context, poolID uuid.UUID, order common.Order, offset, limit *int) ([]common.Chemicals, error)
	DeleteChemicals(ctx context.Context, poolID uuid.UUID, createdAt time.Time) (bool, error)
}

type estimator interface {
	RecommendedChemicals(ctx context.Context, poolID uuid.UUID) (map[common.ChemicalProduct]float64, error)
	DemandMeasurement(ctx context.Context, poolID uuid.UUID) (common.Measurement, error)
	EstimateMeasurement(
		ctx context.Context, poolID uuid.UUID,
		chemicals map[common.ChemicalProduct]float64,
		selection []common.MeasurementType,
	) (common.Measurement, error)
}

type actions interface {
	LogActions(ctx context.Context, poolID uuid.UUID, actions []common.ActionType) (time.Time, error)
	QueryActions(ctx context.Context, poolID uuid.UUID, order common.Order, offset *int, limit *int) ([]common.Action, error)
	DeleteAction(ctx context.Context, poolID uuid.UUID, createdAt time.Time) (bool, error)
}

type poolSettingsManager interface {
	SetSettings(ctx context.Context, poolID uuid.UUID, settings *common.PoolSettings) (*common.PoolSettings, error)
	GetSettings(ctx context.Context, poolID uuid.UUID) (*common.PoolSettings, error)
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
	actions
	poolSettingsManager

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

func (r *Resolver) getMeasurements(ctx context.Context) []common.MeasurementType {
	measurementsMap := map[common.MeasurementType]struct{}{}

	fieldSelections := gql.GetFieldContext(ctx).Field.Selections

	for _, sel := range fieldSelections {
		switch sel := sel.(type) {
		case *ast.Field:
			measurementsMap[mapToMeasurementType(sel.Name)] = struct{}{}
		}
	}

	measurements := make([]common.MeasurementType, 0, len(measurementsMap))
	for m := range measurementsMap {
		measurements = append(measurements, m)
	}

	return measurements
}

func mapToMeasurementType(s string) common.MeasurementType {
	switch s {
	case "chlorine":
		return common.MeasurementChlorine
	case "ph":
		return common.MeasurementPH
	case "alkalinity":
		return common.MeasurementAlkalinity
	default:
		return common.MeasurementChlorine
	}
}

func NewResolver(
	logger logger,
	data poolData,
	measurementHistory measurementHistory,
	additivesHistory additivesHistory,
	estimator estimator,
	actions actions,
	poolSettingsManager poolSettingsManager,
	auth auth,
) *Resolver {
	return &Resolver{
		auth:                auth,
		poolData:            data,
		measurementHistory:  measurementHistory,
		additivesHistory:    additivesHistory,
		estimator:           estimator,
		actions:             actions,
		poolSettingsManager: poolSettingsManager,
		log:                 logger,
	}
}
