package model

import (
	"github.com/PoolHealth/PoolHealthServer/common"
	gqlCommon "github.com/PoolHealth/PoolHealthServer/pkg/api/v1/common"
)

func PoolFromCommon(pool *common.Pool) *Pool {
	if pool == nil {
		return nil
	}

	return &Pool{
		ID:     gqlCommon.ID(pool.ID),
		Name:   pool.Name,
		Volume: pool.Volume,
	}
}
