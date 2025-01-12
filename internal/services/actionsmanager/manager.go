package actionsmanager

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type Manager interface {
	LogActions(ctx context.Context, poolID uuid.UUID, actions []common.ActionType) (*time.Time, error)
	QueryActions(ctx context.Context, poolID uuid.UUID, order common.Order, offset *int, limit *int) ([]common.ActionType, error)
}

type repo interface {
}

type manager struct {
	repo repo

	log log.Logger
}

func (m *manager) LogActions(ctx context.Context, poolID uuid.UUID, actions []common.ActionType) (*time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (m *manager) QueryActions(ctx context.Context, poolID uuid.UUID, order common.Order, offset *int, limit *int) ([]common.ActionType, error) {
	//TODO implement me
	panic("implement me")
}

func NewManager(repo repo, logger log.Logger) Manager {
	return &manager{repo: repo, log: logger}
}
