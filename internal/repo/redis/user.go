package redis

import (
	"context"
	"errors"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/PoolHealth/PoolHealthServer/common"
)

type userRepo interface {
	GetOrCreateUser(ctx context.Context, unique string) (uuid.UUID, error)
}

func (d *db) GetOrCreateUser(ctx context.Context, unique string) (uuid.UUID, error) {
	id, err := d.db.Get(ctx, d.keyBuilder.UserByAppleID(unique)).Result()
	if err == nil {
		return uuid.MustParse(id), nil
	}

	if !errors.Is(err, redis.Nil) {
		return uuid.Nil, err
	}

	uid := uuid.New()
	id = uid.String()

	data, err := json.Marshal(&common.User{ID: uid, AppleID: unique})
	if err != nil {
		return uuid.Nil, err
	}

	if err = d.db.MSet(ctx, map[string]any{
		d.keyBuilder.UserByAppleID(unique): id,
		d.keyBuilder.User(uid):             data,
	}).Err(); err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}
