package estimator

import (
	"github.com/PoolHealth/PoolHealthServer/common"
)

const (
	phLow          = 7.2
	phHigh         = 7.8
	AlkalinityLow  = 80
	AlkalinityHigh = 120
)

var chlorineLowByUsage = map[common.UsageType]float64{
	common.UsageTypePrivate: 1,
	common.UsageTypePublic:  3,
	common.UsageTypeHoliday: 1,
}

var chlorineHighByUsage = map[common.UsageType]float64{
	common.UsageTypePrivate: 3.0,
	common.UsageTypePublic:  5,
	common.UsageTypeHoliday: 5,
}

var chemicalProductCoefficients = map[common.ChemicalProduct]float64{
	common.CalciumHypochlorite65Percent: 1000 / 0.00175,
	common.SodiumHypochlorite12Percent:  1000 / 0.00825,
	common.SodiumHypochlorite14Percent:  1000 / 0.00715,
	common.TCCA90PercentTablets:         200 / 0.0011,
	common.MultiActionTablets:           200 / 0.0011,
	common.TCCA90PercentGranules:        1000 / 0.0011,
	common.Dichlor65Percent:             1000 / 0.0011,

	common.SodiumBicarbonate: 1000 / 0.001675,
}

// Chlorine+((CalciumHypochlorite65Percent*1000)/(volume*0,00175)+(SodiumHypochlorite12Percent*1000)/(volume*0,00825)+(K9*1000)/(volume*0,00715)+(L9*200)/(volume*0,0011)+(M9*200)/(volume*0,0011)+(N9*1000)/(volume*0,0011)+(O9*1000)/(volume*0,0011))
// lastChlorine + sum(el * coefficient) / volume
func CalculateChlorine(
	volume float64,
	lastMeasurement float64,
	additives map[common.ChemicalProduct]float64,
) float64 {

	sum := 0.0

	for k, v := range additives {
		if t := common.ChemicalProductTypes[k]; t != common.Chlorine {
			continue
		}

		sum += v * chemicalProductCoefficients[k]
	}

	return lastMeasurement + sum/volume
}

func CalculateAlkalinity(
	volume float64,
	lastMeasurement float64,
	additives map[common.ChemicalProduct]float64,
) float64 {
	value := 0.0
	if v, ok := additives[common.SodiumBicarbonate]; ok {
		value = v * chemicalProductCoefficients[common.SodiumBicarbonate]
	}

	if v, ok := additives[common.HydrochloricAcid]; ok {
		value -= v * 1000 / 0.001
	}

	if v, ok := additives[common.SodiumBisulphate]; ok {
		value -= v * 1000 / 0.0012
	}

	return lastMeasurement + value/volume
}

func CalculatePH(volume float64,
	lastMeasurement float64,
	additives map[common.ChemicalProduct]float64) float64 {
	//=E8-(((1000*Q8)/($A$2*0,02))+((1000*R8)/($A$2*0,05)))
	value := 0.0

	if v, ok := additives[common.SodiumBisulphate]; ok {
		value += v * 1000 / 0.05
	}

	if v, ok := additives[common.HydrochloricAcid]; ok {
		value += v * 1000 / 0.02
	}

	return lastMeasurement - value/volume
}

func recommendChlorineByTarget(volume, lastMeasurement, target float64) map[common.ChemicalProduct]float64 {
	a := (target - lastMeasurement) * volume

	result := make(map[common.ChemicalProduct]float64)

	for k, v := range common.ChemicalProductTypes {
		if v != common.Chlorine {
			continue
		}

		result[k] = a / chemicalProductCoefficients[k]
	}

	return result
}
