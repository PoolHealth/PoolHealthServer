package measurementhistory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type MeasurementHistory interface {
	CreateMeasurement(ctx context.Context, r common.Measurement) (common.Measurement, error)
	QueryMeasurement(ctx context.Context, poolID uuid.UUID, order common.Order, offset, limit *int) ([]common.Measurement, error)
	DeleteMeasurement(ctx context.Context, poolID uuid.UUID, createdAt time.Time) (bool, error)
}

type repo interface {
	CreateMeasurement(ctx context.Context, rec common.Measurement) error
	QueryMeasurement(ctx context.Context, poolID uuid.UUID, order common.Order) ([]common.Measurement, error)
}

type measurementHistory struct {
	repo

	log log.Logger
}

func (m *measurementHistory) QueryMeasurement(ctx context.Context, poolID uuid.UUID, order common.Order, offset, limit *int) ([]common.Measurement, error) {
	return m.repo.QueryMeasurement(ctx, poolID, order)
}

func (m *measurementHistory) CreateMeasurement(ctx context.Context, r common.Measurement) (common.Measurement, error) {
	r.CreatedAt = time.Now()
	if err := m.repo.CreateMeasurement(ctx, r); err != nil {
		return common.Measurement{}, err
	}

	return r, nil
}

func (m *measurementHistory) DeleteMeasurement(ctx context.Context, poolID uuid.UUID, createdAt time.Time) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewMeasurementHistory(repo repo, log log.Logger) MeasurementHistory {
	return &measurementHistory{repo: repo, log: log}
}
