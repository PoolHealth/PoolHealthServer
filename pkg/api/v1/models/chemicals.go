package model

import (
	"github.com/PoolHealth/PoolHealthServer/common"
)

func ChemicalsFromCommon(a *common.Chemicals) *Chemicals {
	result := &Chemicals{
		Value:     ChemicalValuesFromCommonProduct(a.Products),
		CreatedAt: a.CreatedAt,
	}

	return result
}

func ChemicalValuesFromCommonProduct(products map[common.ChemicalProduct]float64) []ChemicalValue {
	result := make([]ChemicalValue, 0, len(products))

	for k, v := range products {
		switch k {
		case common.CalciumHypochlorite65Percent:
			result = append(result, ChlorineChemicalValue{
				Type:  ChlorineChemicalCalciumHypochlorite65Percent,
				Value: v,
			})
		case common.SodiumHypochlorite12Percent:
			result = append(result, ChlorineChemicalValue{
				Type:  ChlorineChemicalSodiumHypochlorite12Percent,
				Value: v,
			})
		case common.SodiumHypochlorite14Percent:
			result = append(result, ChlorineChemicalValue{
				Type:  ChlorineChemicalSodiumHypochlorite14Percent,
				Value: v,
			})
		case common.TCCA90PercentTablets:
			result = append(result, ChlorineChemicalValue{
				Type:  ChlorineChemicalTCCA90PercentTablets,
				Value: v,
			})
		case common.MultiActionTablets:
			result = append(result, ChlorineChemicalValue{
				Type:  ChlorineChemicalMultiActionTablets,
				Value: v,
			})
		case common.TCCA90PercentGranules:
			result = append(result, ChlorineChemicalValue{
				Type:  ChlorineChemicalTCCA90PercentGranules,
				Value: v,
			})
		case common.Dichlor65Percent:
			result = append(result, ChlorineChemicalValue{
				Type:  ChlorineChemicalDichlor65Percent,
				Value: v,
			})
		case common.HydrochloricAcid:
			result = append(result, AcidChemicalValue{
				Type:  AcidChemicalHydrochloricAcid,
				Value: v,
			})
		case common.SodiumBisulphate:
			result = append(result, AcidChemicalValue{
				Type:  AcidChemicalSodiumBisulphate,
				Value: v,
			})
		case common.SodiumBicarbonate:
			result = append(result, AlkalinityChemicalValue{
				Type:  AlkalinityChemicalSodiumBicarbonate,
				Value: v,
			})
		}
	}

	return result
}

func (a *Chemicals) ToCommon() *common.Chemicals {
	result := &common.Chemicals{
		CreatedAt: a.CreatedAt,
		Products:  make(map[common.ChemicalProduct]float64),
	}

	for _, v := range a.Value {
		switch ch := v.(type) {
		case ChlorineChemicalValue:
			switch ch.Type {
			case ChlorineChemicalCalciumHypochlorite65Percent:
				result.Products[common.CalciumHypochlorite65Percent] = ch.Value
			case ChlorineChemicalSodiumHypochlorite12Percent:
				result.Products[common.SodiumHypochlorite12Percent] = ch.Value
			case ChlorineChemicalSodiumHypochlorite14Percent:
				result.Products[common.SodiumHypochlorite14Percent] = ch.Value
			case ChlorineChemicalTCCA90PercentTablets:
				result.Products[common.TCCA90PercentTablets] = ch.Value
			case ChlorineChemicalMultiActionTablets:
				result.Products[common.MultiActionTablets] = ch.Value
			case ChlorineChemicalTCCA90PercentGranules:
				result.Products[common.TCCA90PercentGranules] = ch.Value
			case ChlorineChemicalDichlor65Percent:
				result.Products[common.Dichlor65Percent] = ch.Value
			}
		case AcidChemicalValue:
			switch ch.Type {
			case AcidChemicalHydrochloricAcid:
				result.Products[common.HydrochloricAcid] = ch.Value
			case AcidChemicalSodiumBisulphate:
				result.Products[common.SodiumBisulphate] = ch.Value
			}
		case AlkalinityChemicalValue:
			switch ch.Type { //nolint:gocritic
			case AlkalinityChemicalSodiumBicarbonate:
				result.Products[common.SodiumBicarbonate] = ch.Value
			}
		}
	}

	return result
}

func (i *ChemicalInput) ToCommonProducts() map[common.ChemicalProduct]float64 {
	result := make(map[common.ChemicalProduct]float64)

	for _, v := range i.Chlorine {
		switch v.Type {
		case ChlorineChemicalCalciumHypochlorite65Percent:
			result[common.CalciumHypochlorite65Percent] = v.Value
		case ChlorineChemicalSodiumHypochlorite12Percent:
			result[common.SodiumHypochlorite12Percent] = v.Value
		case ChlorineChemicalSodiumHypochlorite14Percent:
			result[common.SodiumHypochlorite14Percent] = v.Value
		case ChlorineChemicalTCCA90PercentTablets:
			result[common.TCCA90PercentTablets] = v.Value
		case ChlorineChemicalMultiActionTablets:
			result[common.MultiActionTablets] = v.Value
		case ChlorineChemicalTCCA90PercentGranules:
			result[common.TCCA90PercentGranules] = v.Value
		case ChlorineChemicalDichlor65Percent:
			result[common.Dichlor65Percent] = v.Value
		}
	}

	for _, v := range i.Acid {
		switch v.Type {
		case AcidChemicalHydrochloricAcid:
			result[common.HydrochloricAcid] = v.Value
		case AcidChemicalSodiumBisulphate:
			result[common.SodiumBisulphate] = v.Value
		}
	}

	for _, v := range i.Alkalinity {
		switch v.Type { //nolint:gocritic
		case AlkalinityChemicalSodiumBicarbonate:
			result[common.SodiumBicarbonate] = v.Value
		}
	}

	return result
}
