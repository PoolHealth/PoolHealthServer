package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/PoolHealth/PoolHealthServer/internal/repo/keys"
	"github.com/PoolHealth/PoolHealthServer/internal/repo/redis/migrations"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

const versionKey = "version"

type DB interface {
	Migrate(ctx context.Context) error
	version
	poolRepo
	userRepo
	poolSettings
}

type db struct {
	keyBuilder keys.Builder
	db         *redis.Client

	log log.Logger
}

func (d *db) Migrate(ctx context.Context) error {
	v, err := d.GetVersion(ctx)
	if err != nil {
		return err
	}

	d.log.WithField(versionKey, v).Info("current version")

	m := migrations.Migrations(v)
	for _, el := range m {
		d.log.WithField(versionKey, el.Version()).Info("migrating to version")

		tx := d.db.TxPipeline()

		if err = el.Up(ctx, tx); err != nil {
			return err
		}

		if err = d.WriteVersion(ctx, el.Version()); err != nil {
			return err
		}

		if _, err = tx.Exec(ctx); err != nil {
			return err
		}

		d.log.WithField(versionKey, el.Version()).Info("migrated to version")
	}

	return nil
}

func NewDB(client *redis.Client, log log.Logger) DB {
	return &db{
		keyBuilder: keys.NewBuilder(),
		db:         client,
		log:        log,
	}
}
