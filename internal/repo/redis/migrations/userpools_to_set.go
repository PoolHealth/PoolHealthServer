package migrations

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"

	keybuilder "github.com/PoolHealth/PoolHealthServer/internal/repo/keys"
)

type UserPoolToSet struct {
}

func (u UserPoolToSet) Up(ctx context.Context, tr *redis.Tx) error {
	kb := keybuilder.NewBuilder()

	var cursor uint64

	for {
		var (
			keys []string
			err  error
		)

		keys, cursor, err = tr.Scan(ctx, cursor, kb.UsersPools()+"*", 0).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			userPools := make([]string, 0)

			data, err := tr.Get(ctx, key).Bytes()

			if data != nil {
				if err = json.Unmarshal(data, &userPools); err != nil {
					return err
				}

				values := make([]any, len(userPools))
				for i, v := range userPools {
					values[i] = v
				}

				if err := tr.Del(ctx, key).Err(); err != nil {
					return err
				}

				if err = tr.SAdd(ctx, key, values...).Err(); err != nil {
					return err
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (u UserPoolToSet) Down(context.Context, *redis.Tx) error {
	return nil
}

func (u UserPoolToSet) Version() uint32 {
	return 1
}
