package estimator

import (
	"errors"

	"github.com/PoolHealth/PoolHealthServer/common"
)

// Chlorine+((CalciumHypochlorite65Percent*1000)/(volume*0,00175)+(SodiumHypochlorite12Percent*1000)/(volume*0,00825)+(K9*1000)/(volume*0,00715)+(L9*200)/(volume*0,0011)+(M9*200)/(volume*0,0011)+(N9*1000)/(volume*0,0011)+(O9*1000)/(volume*0,0011))
// lastChlorine + sum(el * coefficient) / volume
func CalculateChlorine(volume float64, lastMeasurement common.Measurement, lastAdditives common.Additives) (float64, error) {
	if !lastMeasurement.Chlorine.Valid {
		return 0.0, errors.New("last chlorine measurement required")
	}

	sum := 0.0

	for k, v := range lastAdditives.Products {
		sum += v * common.ChemicalProductCoefficients[k]
	}

	return lastMeasurement.Chlorine.Float64 + sum/volume, nil
}
