package redis

import (
	"context"

	"github.com/PoolHealth/PoolHealthServer/common"
)

type measurements interface {
	CreateMeasurement(ctx context.Context, rec *common.Measurement) error
	ListMeasurements(ctx context.Context, poolID string) ([]common.Measurement, error)
}
