package sheetsmigrator

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/internal/models"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type Migrator struct {
	sheetAdapter sheetAdapter
	pool         poolManager
	history      historyRepo
	migrations   map[uuid.UUID]map[uuid.UUID]common.Migration
	mu           sync.RWMutex

	lg log.Logger
}

type sheetAdapter interface {
	GetPools(ctx context.Context, sheetID string) ([]models.Pool, error)
}

type poolManager interface {
	Create(ctx context.Context, userID uuid.UUID, data *common.PoolData) (pool *common.Pool, err error)
	Update(ctx context.Context, id uuid.UUID, rec *common.PoolData) (record *common.Pool, err error)
	SearchPoolByName(ctx context.Context, userID uuid.UUID, name string) (*common.Pool, error)
}

type historyRepo interface {
	CreateMeasurement(ctx context.Context, rec common.Measurement) error
	CreateChemicals(ctx context.Context, rec *common.Chemicals) error
	LogActions(ctx context.Context, poolID uuid.UUID, actions *common.Action) error
}

func NewMigrator(
	sheetAdapter sheetAdapter,
	pool poolManager,
	history historyRepo,

	lg log.Logger,
) *Migrator {
	return &Migrator{
		sheetAdapter: sheetAdapter,
		pool:         pool,
		history:      history,
		migrations:   map[uuid.UUID]map[uuid.UUID]common.Migration{},
		lg:           lg,
	}
}

func (m *Migrator) Migrate(_ context.Context, userID uuid.UUID, sheetID string) uuid.UUID {
	id := uuid.New()

	m.mu.Lock()
	if _, ok := m.migrations[userID]; !ok {
		m.migrations[userID] = map[uuid.UUID]common.Migration{
			id: {
				ID:     id,
				Status: common.MigrationStatusPending,
			},
		}
	}

	m.migrations[userID][id] = common.Migration{
		ID:     id,
		Status: common.MigrationStatusPending,
	}

	go m.migrate(id, userID, sheetID)
	m.mu.Unlock()

	return id
}

func (m *Migrator) Migration(ctx context.Context, userID, id uuid.UUID) (common.Migration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if migrations, ok := m.migrations[userID]; ok {
		if migration, ok := migrations[id]; ok {
			return migration, nil
		}
	}

	return common.Migration{ID: id, Status: common.MigrationStatusUnknown}, nil
}

func (m *Migrator) migrate(id uuid.UUID, userID uuid.UUID, sheetID string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	pools, err := m.sheetAdapter.GetPools(ctx, sheetID)
	if err != nil {
		m.changeStatus(id, userID, common.MigrationStatusFailed)
		m.lg.WithError(err).Error("Get pools from sheet")

		return
	}

	for _, pool := range pools {
		p, err := m.pool.SearchPoolByName(ctx, userID, pool.Pool.Name)
		if err != nil {
			if !errors.Is(err, common.ErrPoolNotFound) {
				m.changeStatus(id, userID, common.MigrationStatusFailed)
				m.lg.WithError(err).Error("Can't find pool")

				return
			}

			p, err = m.pool.Create(ctx, userID, &pool.Pool)
			if err != nil {
				m.changeStatus(id, userID, common.MigrationStatusFailed)
				m.lg.WithError(err).Error("Can't create pool")

				return
			}
		} else {
			_, err = m.pool.Update(ctx, p.ID, &pool.Pool)
			if err != nil {
				m.changeStatus(id, userID, common.MigrationStatusFailed)
				m.lg.WithError(err).Error("Can't update pool")

				return
			}
		}

		for _, measurement := range pool.Measurements {
			measurement.PoolID = p.ID

			if err = m.history.CreateMeasurement(ctx, measurement); err != nil {
				if errors.Is(err, common.ErrRecordAlreadyExists) {
					continue
				}

				m.changeStatus(id, userID, common.MigrationStatusFailed)
				m.lg.WithError(err).Error("Can't create measurement")

				return
			}
		}

		for _, chemical := range pool.Chemicals {
			chemical.PoolID = p.ID

			if err = m.history.CreateChemicals(ctx, &chemical); err != nil {
				if errors.Is(err, common.ErrRecordAlreadyExists) {
					continue
				}

				m.changeStatus(id, userID, common.MigrationStatusFailed)
				m.lg.WithError(err).Error("Can't create chemicals")

				return
			}
		}

		for _, action := range pool.Actions {
			if err := m.history.LogActions(ctx, p.ID, &action); err != nil {
				if errors.Is(err, common.ErrRecordAlreadyExists) {
					continue
				}

				m.changeStatus(id, userID, common.MigrationStatusFailed)
				m.lg.WithError(err).Error("Can't log actions")

				return
			}
		}
	}

	m.changeStatus(id, userID, common.MigrationsStatusCompleted)
}

func (m *Migrator) changeStatus(id uuid.UUID, userID uuid.UUID, status common.MigrationStatus) {
	m.mu.Lock()
	m.migrations[userID][id] = common.Migration{
		ID:     id,
		Status: status,
	}
	m.mu.Unlock()
}
