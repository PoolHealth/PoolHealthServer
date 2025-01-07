package model

import (
	"github.com/PoolHealth/PoolHealthServer/common"
)

func MeasurementFromCommon(res common.Measurement) *Measurement {
	return &Measurement{
		Chlorine:   res.Chlorine.Float64,
		Ph:         res.PH.Float64,
		Alkalinity: res.Alkalinity.Float64,
		CreatedAt:  res.CreatedAt,
	}
}
