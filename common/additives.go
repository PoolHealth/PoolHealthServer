package common

import (
	"time"

	"github.com/google/uuid"
)

type Additives struct {
	PoolID    uuid.UUID
	Products  map[ChemicalProduct]float64
	CreatedAt time.Time
}

type ChemicalProduct uint64

const (
	CalciumHypochlorite65Percent ChemicalProduct = iota
	SodiumHypochlorite12Percent
	SodiumHypochlorite14Percent
	TCCA90PercentTablets
	MultiActionTablets
	TCCA90PercentGranules
	Dichlor65Percent
)

var ChemicalProductNames = map[ChemicalProduct]string{
	CalciumHypochlorite65Percent: "Calcium Hypochlorite 65%",
	SodiumHypochlorite12Percent:  "Sodium Hypochlorite 12%",
	SodiumHypochlorite14Percent:  "Sodium Hypochlorite 14%",
	TCCA90PercentTablets:         "TCCA 90% Tablets",
	MultiActionTablets:           "Multi-Action Tablets",
	TCCA90PercentGranules:        "TCCA 90% Granules",
	Dichlor65Percent:             "Dichlor 65%",
}

var ChemicalProductCoefficients = map[ChemicalProduct]float64{
	CalciumHypochlorite65Percent: 1000 / 0.00175,
	SodiumHypochlorite12Percent:  1000 / 0.00825,
	SodiumHypochlorite14Percent:  1000 / 0.00715,
	TCCA90PercentTablets:         200 / 0.0011,
	MultiActionTablets:           200 / 0.0011,
	TCCA90PercentGranules:        1000 / 0.0011,
	Dichlor65Percent:             1000 / 0.0011,
}

var ChemicalProductNamesToChemicalProduct = map[string]ChemicalProduct{
	"Calcium Hypochlorite 65%": CalciumHypochlorite65Percent,
	"Sodium Hypochlorite 12%":  SodiumHypochlorite12Percent,
	"Sodium Hypochlorite 14%":  SodiumHypochlorite14Percent,
	"TCCA 90% Tablets":         TCCA90PercentTablets,
	"Multi-Action Tablets":     MultiActionTablets,
	"TCCA 90% Granules":        TCCA90PercentGranules,
	"Dichlor 65%":              Dichlor65Percent,
}
