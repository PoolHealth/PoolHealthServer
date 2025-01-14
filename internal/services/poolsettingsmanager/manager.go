package poolsettingsmanager

import (
	"context"

	"github.com/google/uuid"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type Manager struct {
	repo repo

	log log.Logger
}

type repo interface {
	SetSettings(ctx context.Context, poolID uuid.UUID, settings *common.PoolSettings) error
	GetSettings(ctx context.Context, poolID uuid.UUID) (*common.PoolSettings, error)
}

func NewPoolSettingsManager(repo repo, logger log.Logger) *Manager {
	return &Manager{repo: repo, log: logger}
}

func (m *Manager) SetSettings(ctx context.Context, poolID uuid.UUID, settings *common.PoolSettings) (*common.PoolSettings, error) {
	return settings, m.repo.SetSettings(ctx, poolID, settings)
}

func (m *Manager) GetSettings(ctx context.Context, poolID uuid.UUID) (*common.PoolSettings, error) {
	return m.repo.GetSettings(ctx, poolID)
}
