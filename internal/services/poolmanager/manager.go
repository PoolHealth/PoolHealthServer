package poolmanager

import (
	"context"

	"github.com/google/uuid"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

type Manager interface {
	Create(ctx context.Context, userID uuid.UUID, data *common.PoolData) (pool *common.Pool, err error)
	Update(ctx context.Context, id uuid.UUID, rec *common.PoolData) (record *common.Pool, err error)
	Delete(ctx context.Context, id uuid.UUID) error
	Has(ctx context.Context, id uuid.UUID, userID uuid.UUID) (bool, error)
	List(ctx context.Context, userID uuid.UUID) ([]common.Pool, error)
	SearchPoolByName(ctx context.Context, userID uuid.UUID, name string) (*common.Pool, error)
	SubscribeOnCreate(ctx context.Context) (<-chan *common.Pool, error)
	SubscribeOnUpdate(ctx context.Context) (<-chan *common.Pool, error)
	SubscribeOnDelete(ctx context.Context) (<-chan uuid.UUID, error)
}

type repo interface {
	CreatePool(ctx context.Context, userID uuid.UUID, id uuid.UUID, rec *common.PoolData) error
	UpdatePool(ctx context.Context, id uuid.UUID, rec *common.PoolData) error
	DeletePool(ctx context.Context, id uuid.UUID) error
	ListPool(ctx context.Context, userID uuid.UUID) ([]common.Pool, error)
	UserHasPool(ctx context.Context, id uuid.UUID, userID uuid.UUID) (ok bool, err error)
}

type manager struct {
	repo

	log log.Logger
}

func NewManager(repo repo, logger log.Logger) Manager {
	return &manager{repo: repo, log: logger}
}

func (m *manager) Create(ctx context.Context, userID uuid.UUID, data *common.PoolData) (pool *common.Pool, err error) {
	id := uuid.New()
	if err = m.repo.CreatePool(ctx, id, userID, data); err != nil {
		return nil, err
	}

	return &common.Pool{
		ID:       id,
		PoolData: *data,
	}, nil
}

func (m *manager) Update(ctx context.Context, id uuid.UUID, rec *common.PoolData) (record *common.Pool, err error) {
	if err = m.repo.UpdatePool(ctx, id, rec); err != nil {
		return nil, err
	}

	return &common.Pool{
		ID:       id,
		PoolData: *rec,
	}, nil
}

func (m *manager) Delete(ctx context.Context, id uuid.UUID) error {
	return m.repo.DeletePool(ctx, id)
}

func (m *manager) Has(ctx context.Context, id uuid.UUID, userID uuid.UUID) (bool, error) {
	return m.repo.UserHasPool(ctx, id, userID)
}

func (m *manager) List(ctx context.Context, userID uuid.UUID) ([]common.Pool, error) {
	return m.repo.ListPool(ctx, userID)
}

func (m *manager) SearchPoolByName(ctx context.Context, userID uuid.UUID, name string) (*common.Pool, error) {
	pools, err := m.repo.ListPool(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, pool := range pools {
		if pool.Name == name {
			return &pool, nil
		}
	}

	return nil, common.ErrPoolNotFound
}

func (m *manager) SubscribeOnCreate(ctx context.Context) (<-chan *common.Pool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *manager) SubscribeOnUpdate(ctx context.Context) (<-chan *common.Pool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *manager) SubscribeOnDelete(ctx context.Context) (<-chan uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}
