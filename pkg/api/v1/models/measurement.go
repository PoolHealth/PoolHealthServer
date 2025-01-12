package model

import (
	"github.com/PoolHealth/PoolHealthServer/common"
)

func MeasurementRecordFromCommon(res common.Measurement) *MeasurementRecord {
	return &MeasurementRecord{
		Measurement: MeasurementFromCommon(res),
		CreatedAt:   res.CreatedAt,
	}
}

func MeasurementFromCommon(res common.Measurement) *Measurement {
	return &Measurement{
		Chlorine:   res.Chlorine.Float64,
		Ph:         res.PH.Float64,
		Alkalinity: res.Alkalinity.Float64,
	}
}
