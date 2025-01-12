package common

import "github.com/google/uuid"

type Pool struct {
	ID uuid.UUID
	PoolMetadata
	PoolData
}

type PoolData struct {
	Name   string
	Volume float64
}

type (
	PoolMetadata struct {
		CleanerUserID uuid.UUID
	}
	PoolSettings struct {
		Type         PoolType
		UsageType    UsageType
		LocationType LocationType
		PoolShape    PoolShape
		Coordinates  *CoordinatesInput
	}

	CoordinatesInput struct {
		Latitude  float64
		Longitude float64
	}
)

type PoolType uint

const (
	// PoolType enum
	PoolTypeUnknown PoolType = iota
	PoolTypeInfinity
	PoolTypeOverflow
	PoolTypeSkimmer
)

type UsageType uint

const (
	// UsageType enum
	UsageTypeUnknown UsageType = iota
	UsageTypePrivate
	UsageTypePublic
	UsageTypeHoliday
)

type LocationType uint

const (
	// LocationType enum
	LocationTypeUnknown LocationType = iota
	LocationTypeIndoor
	LocationTypeOutdoor
)

type PoolShape uint

const (
	// PoolShape enum
	PoolShapeUnknown PoolShape = iota
	PoolShapeRectangle
	PoolShapeCircle
	PoolShapeOval
	PoolShapeKidney
	PoolShapeLShaped
	PoolShapeTShaped
	PoolShapeFreeForm
)
