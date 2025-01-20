package actionsmanager

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type Manager interface {
	DeleteAction(ctx context.Context, poolID uuid.UUID, createdAt time.Time) (bool, error)
	LogActions(ctx context.Context, poolID uuid.UUID, actions []common.ActionType) (time.Time, error)
	QueryActions(ctx context.Context, poolID uuid.UUID, order common.Order, offset, limit *int) ([]common.Action, error)
}

type repo interface {
	LogActions(ctx context.Context, poolID uuid.UUID, actions *common.Action) error
	QueryActions(
		ctx context.Context,
		poolID uuid.UUID,
		order common.Order,
	) ([]common.Action, error)
	DeleteAction(ctx context.Context, poolID uuid.UUID, createdAt time.Time) error
}

type manager struct {
	repo repo

	log log.Logger
}

func (m *manager) DeleteAction(ctx context.Context, poolID uuid.UUID, createdAt time.Time) (bool, error) {
	return true, m.repo.DeleteAction(ctx, poolID, createdAt)
}

func (m *manager) LogActions(ctx context.Context, poolID uuid.UUID, actions []common.ActionType) (time.Time, error) {
	action := &common.Action{
		Types:     actions,
		CreatedAt: time.Now(),
	}

	if err := m.repo.LogActions(ctx, poolID, action); err != nil {
		return time.Time{}, err
	}

	return action.CreatedAt, nil
}

func (m *manager) QueryActions(
	ctx context.Context,
	poolID uuid.UUID,
	order common.Order,
	_, _ *int,
) ([]common.Action, error) {
	return m.repo.QueryActions(ctx, poolID, order)
}

func NewManager(repo repo, logger log.Logger) Manager {
	return &manager{repo: repo, log: logger}
}
