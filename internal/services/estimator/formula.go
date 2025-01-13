package estimator

import (
	"github.com/PoolHealth/PoolHealthServer/common"
)

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

		sum += v * common.ChemicalProductCoefficients[k]
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
		value = v * common.ChemicalProductCoefficients[common.SodiumBicarbonate]
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
