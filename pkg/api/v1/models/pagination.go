package model

import "github.com/PoolHealth/PoolHealthServer/common"

func (o Order) ToCommon() common.Order {
	switch o {
	case OrderAsc:
		return common.OrderAsc
	case OrderDesc:
		return common.OrderDesc
	default:
		return common.OrderAsc
	}
}
