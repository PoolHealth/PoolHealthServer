package estimator

import (
	"errors"

	"github.com/PoolHealth/PoolHealthServer/common"
)

// Chlorine+((CalciumHypochlorite65Percent*1000)/(volume*0,00175)+(SodiumHypochlorite12Percent*1000)/(volume*0,00825)+(K9*1000)/(volume*0,00715)+(L9*200)/(volume*0,0011)+(M9*200)/(volume*0,0011)+(N9*1000)/(volume*0,0011)+(O9*1000)/(volume*0,0011))
// lastChlorine + sum(el * coefficient) / volume
func CalculateChlorine(
	volume float64,
	lastMeasurement common.Measurement,
	additives map[common.ChemicalProduct]float64,
) (float64, error) {
	if !lastMeasurement.Chlorine.Valid {
		return 0.0, errors.New("last chlorine measurement required")
	}

	sum := 0.0

	for k, v := range additives {
		if t := common.ChemicalProductTypes[k]; t != common.Chlorine {
			continue
		}

		sum += v * common.ChemicalProductCoefficients[k]
	}

	return lastMeasurement.Chlorine.Float64 + sum/volume, nil
}

func CalculateAlkalinity(
	volume float64,
	lastMeasurement common.Measurement,
	additives map[common.ChemicalProduct]float64,
) (float64, error) {
	if !lastMeasurement.Alkalinity.Valid {
		return 0.0, errors.New("last alkalinity measurement required")
	}

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

	return lastMeasurement.Alkalinity.Float64 + value/volume, nil
}

func CalculatePH(volume float64,
	lastMeasurement common.Measurement,
	additives map[common.ChemicalProduct]float64) (float64, error) {
	if !lastMeasurement.PH.Valid {
		return 0.0, errors.New("last PH measurement required")
	}

	//=E8-(((1000*Q8)/($A$2*0,02))+((1000*R8)/($A$2*0,05)))
	value := 0.0

	if v, ok := additives[common.SodiumBisulphate]; ok {
		value += v * 1000 / 0.05
	}

	if v, ok := additives[common.HydrochloricAcid]; ok {
		value += v * 1000 / 0.02
	}

	return lastMeasurement.PH.Float64 - value/volume, nil
}
