package models

import "github.com/PoolHealth/PoolHealthServer/common"

type Pool struct {
	Pool         common.PoolData
	Measurements []common.Measurement
	Chemicals    []common.Chemicals
	Actions      []common.Action
}
