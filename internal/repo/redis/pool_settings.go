package redis

import (
	"context"
	"errors"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/PoolHealth/PoolHealthServer/common"
)

type poolSettings interface {
	SetSettings(ctx context.Context, poolID uuid.UUID, settings *common.PoolSettings) error
	GetSettings(ctx context.Context, poolID uuid.UUID) (*common.PoolSettings, error)
}

func (d *db) SetSettings(ctx context.Context, poolID uuid.UUID, settings *common.PoolSettings) error {
	key := d.keyBuilder.PoolSettings(poolID)

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	return d.db.Set(ctx, key, data, 0).Err()
}

func (d *db) GetSettings(ctx context.Context, poolID uuid.UUID) (*common.PoolSettings, error) {
	key := d.keyBuilder.PoolSettings(poolID)

	data, err := d.db.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var settings *common.PoolSettings
	if err = json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return settings, nil
}
