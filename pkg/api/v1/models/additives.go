package model

import (
	"github.com/PoolHealth/PoolHealthServer/common"
)

func AdditivesFromCommon(a *common.Additives) *Additives {
	result := &Additives{
		CreatedAt: a.CreatedAt,
	}

	for k, v := range a.Products {
		switch k {
		case common.CalciumHypochlorite65Percent:
			result.CalciumHypochlorite65Percent = &v
		case common.SodiumHypochlorite12Percent:
			result.SodiumHypochlorite12Percent = &v
		case common.SodiumHypochlorite14Percent:
			result.SodiumHypochlorite14Percent = &v
		case common.TCCA90PercentTablets:
			result.TCCA90PercentTablets = &v
		case common.MultiActionTablets:
			result.MultiActionTablets = &v
		case common.TCCA90PercentGranules:
			result.TCCA90PercentGranules = &v
		case common.Dichlor65Percent:
			result.Dichlor65Percent = &v
		}
	}

	return result
}

func (a *Additives) ToCommon() *common.Additives {
	result := &common.Additives{
		CreatedAt: a.CreatedAt,
		Products:  make(map[common.ChemicalProduct]float64),
	}

	if a.CalciumHypochlorite65Percent != nil {
		result.Products[common.CalciumHypochlorite65Percent] = *a.CalciumHypochlorite65Percent
	}
	if a.SodiumHypochlorite12Percent != nil {
		result.Products[common.SodiumHypochlorite12Percent] = *a.SodiumHypochlorite12Percent
	}
	if a.SodiumHypochlorite14Percent != nil {
		result.Products[common.SodiumHypochlorite14Percent] = *a.SodiumHypochlorite14Percent
	}
	if a.TCCA90PercentTablets != nil {
		result.Products[common.TCCA90PercentTablets] = *a.TCCA90PercentTablets
	}
	if a.MultiActionTablets != nil {
		result.Products[common.MultiActionTablets] = *a.MultiActionTablets
	}
	if a.TCCA90PercentGranules != nil {
		result.Products[common.TCCA90PercentGranules] = *a.TCCA90PercentGranules
	}
	if a.Dichlor65Percent != nil {
		result.Products[common.Dichlor65Percent] = *a.Dichlor65Percent
	}

	return result
}
