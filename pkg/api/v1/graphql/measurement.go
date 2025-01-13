package graphql

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
	m := &Measurement{}

	if res.Chlorine.Valid {
		m.Chlorine = &res.Chlorine.Float64
	}

	if res.PH.Valid {
		m.Ph = &res.PH.Float64
	}

	if res.Alkalinity.Valid {
		m.Alkalinity = &res.Alkalinity.Float64
	}

	return m
}
