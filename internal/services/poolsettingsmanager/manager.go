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
}

func NewPoolSettingsManager(repo repo, logger log.Logger) *Manager {
	return &Manager{repo: repo, log: logger}
}

func (m *Manager) SetSettings(ctx context.Context, poolID uuid.UUID, settings *common.PoolSettings) (*common.PoolSettings, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Manager) GetSettings(ctx context.Context, poolID uuid.UUID) (*common.PoolSettings, error) {
	//TODO implement me
	panic("implement me")
}
