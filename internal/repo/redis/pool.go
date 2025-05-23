package redis

import (
	"context"
	"errors"
	"strings"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/PoolHealth/PoolHealthServer/common"
)

type poolRepo interface {
	CreatePool(ctx context.Context, id uuid.UUID, userID uuid.UUID, rec *common.PoolData) error
	UpdatePool(ctx context.Context, id uuid.UUID, rec *common.PoolData) error
	SetPool(ctx context.Context, id uuid.UUID, rec *common.PoolData) error
	GetPool(ctx context.Context, id uuid.UUID) (*common.Pool, error)
	HasPool(ctx context.Context, id uuid.UUID) (ok bool, err error)
	UserHasPool(ctx context.Context, id uuid.UUID, userID uuid.UUID) (ok bool, err error)
	DeletePool(ctx context.Context, id, userID uuid.UUID) error
	ListPool(ctx context.Context, userID uuid.UUID) ([]common.Pool, error)
}

func (d *db) CreatePool(ctx context.Context, id uuid.UUID, userID uuid.UUID, rec *common.PoolData) error {
	key := d.keyBuilder.UserPools(userID)

	return d.db.Watch(ctx, func(tx *redis.Tx) error {
		ok, err := d.HasPool(ctx, id)
		if err != nil {
			return err
		}

		if ok {
			return common.ErrPoolAlreadyExists
		}

		if err = d.SetPool(ctx, id, rec); err != nil {
			return err
		}

		return tx.SAdd(ctx, key, d.keyBuilder.Pool(id)).Err()
	}, key)
}

func (d *db) UpdatePool(ctx context.Context, id uuid.UUID, rec *common.PoolData) error {
	key := d.keyBuilder.Pool(id)
	return d.db.Watch(ctx, func(tx *redis.Tx) error {

		ok, err := d.hasPool(ctx, tx, id)
		if err != nil {
			return err
		}

		if !ok {
			return common.ErrPoolNotFound
		}

		return d.setPool(ctx, tx, id, rec)
	}, key)
}

func (d *db) SetPool(ctx context.Context, id uuid.UUID, rec *common.PoolData) error {
	return d.setPool(ctx, d.db, id, rec)
}

func (d *db) setPool(ctx context.Context, db redis.StringCmdable, id uuid.UUID, rec *common.PoolData) error {
	key := d.keyBuilder.Pool(id)
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	d.log.WithField("id", id).Trace(string(data))

	return db.Set(ctx, key, data, 0).Err()
}

func (d *db) GetPool(ctx context.Context, id uuid.UUID) (record *common.Pool, err error) {
	key := d.keyBuilder.Pool(id)
	data, err := d.db.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, common.ErrPoolNotFound
	}

	var rec common.PoolData
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}

	return &common.Pool{
		ID:       id,
		PoolData: rec,
	}, nil
}

func (d *db) hasPool(ctx context.Context, db redis.StringCmdable, id uuid.UUID) (ok bool, err error) {
	data, err := db.Get(ctx, d.keyBuilder.Pool(id)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}

	return data != nil, nil
}

func (d *db) HasPool(ctx context.Context, id uuid.UUID) (ok bool, err error) {
	return d.hasPool(ctx, d.db, id)
}

func (d *db) DeletePool(ctx context.Context, id, userID uuid.UUID) error {
	userKey := d.keyBuilder.UserPools(userID)
	key := d.keyBuilder.Pool(id)

	return d.db.Watch(ctx, func(tx *redis.Tx) error {
		if err := d.db.Del(ctx, d.keyBuilder.Pool(id)).Err(); err != nil {
			return err
		}

		return tx.SPop(ctx, userKey).Err()
	}, key, userKey)
}

func (d *db) UserHasPool(ctx context.Context, id uuid.UUID, userID uuid.UUID) (ok bool, err error) {
	key := d.keyBuilder.UserPools(userID)
	exist, err := d.db.SIsMember(ctx, key, d.keyBuilder.Pool(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}

	return exist, nil
}

func (d *db) ListPool(ctx context.Context, userID uuid.UUID) ([]common.Pool, error) {
	userKey := d.keyBuilder.UserPools(userID)

	poolKeys, err := d.db.SMembers(ctx, userKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []common.Pool{}, nil
		}

		return nil, err
	}

	prefix := d.keyBuilder.Pools()

	pools := make([]common.Pool, 0, len(poolKeys))

	for _, key := range poolKeys {
		if len(key) < 4 {
			return nil, common.ErrInvalidPoolID
		}

		id, err := uuid.Parse(strings.Replace(key, prefix, "", 1))
		if err != nil {
			return nil, err
		}

		pool := common.PoolData{}

		data, err := d.db.Get(ctx, key).Bytes()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}

			return nil, err
		}

		if err = json.Unmarshal(data, &pool); err != nil {
			return nil, err
		}

		pools = append(pools, common.Pool{
			ID:           id,
			PoolMetadata: common.PoolMetadata{CleanerUserID: userID},
			PoolData:     pool,
		})
	}

	return pools, nil
}
