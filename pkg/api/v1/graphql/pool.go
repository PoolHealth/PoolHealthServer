package graphql

import (
	"errors"

	"github.com/PoolHealth/PoolHealthServer/common"
	gqlCommon "github.com/PoolHealth/PoolHealthServer/pkg/api/v1/common"
)

var (
	ErrUnknownPoolType     = errors.New("unknown pool type")
	ErrUnknownUsageType    = errors.New("unknown usage type")
	ErrUnknownLocationType = errors.New("unknown location type")
	ErrUnknownPoolShape    = errors.New("unknown pool shape")
)

func PoolFromCommon(pool *common.Pool) *Pool {
	if pool == nil {
		return nil
	}

	return &Pool{
		ID:     gqlCommon.ID(pool.ID),
		Name:   pool.Name,
		Volume: pool.Volume,
	}
}

func (i *PoolSettingsInput) ToCommon() *common.PoolSettings {
	if i == nil {
		return nil
	}

	return &common.PoolSettings{
		Type:         i.Type.ToCommon(),
		UsageType:    i.UsageType.ToCommon(),
		LocationType: i.LocationType.ToCommon(),
		PoolShape:    i.PoolShape.ToCommon(),
		Coordinates:  i.Coordinates.ToCommon(),
	}
}

func (e PoolType) ToCommon() common.PoolType {
	switch e {
	case PoolTypeInfinity:
		return common.PoolTypeInfinity
	case PoolTypeOverflow:
		return common.PoolTypeOverflow
	case PoolTypeSkimmer:
		return common.PoolTypeSkimmer
	}

	return common.PoolTypeUnknown
}

func (e UsageType) ToCommon() common.UsageType {
	switch e {
	case UsageTypePrivate:
		return common.UsageTypePrivate
	case UsageTypeCommunity:
		return common.UsageTypePublic
	case UsageTypeHoliday:
		return common.UsageTypeHoliday
	}

	return common.UsageTypeUnknown
}

func (e LocationType) ToCommon() common.LocationType {
	switch e {
	case LocationTypeIndoor:
		return common.LocationTypeIndoor
	case LocationTypeOutdoor:
		return common.LocationTypeOutdoor
	}

	return common.LocationTypeUnknown
}

func (e PoolShape) ToCommon() common.PoolShape {
	switch e {
	case PoolShapeRectangle:
		return common.PoolShapeRectangle
	case PoolShapeOval:
		return common.PoolShapeOval
	case PoolShapeCircle:
		return common.PoolShapeCircle
	case PoolShapeKidney:
		return common.PoolShapeKidney
	case PoolShapeL:
		return common.PoolShapeLShaped
	case PoolShapeT:
		return common.PoolShapeTShaped
	case PoolShapeFreeForm:
		return common.PoolShapeFreeForm
	}

	return common.PoolShapeUnknown
}

func (i CoordinatesInput) ToCommon() *common.CoordinatesInput {
	if i == (CoordinatesInput{}) {
		return nil
	}

	return &common.CoordinatesInput{
		Latitude:  i.Latitude,
		Longitude: i.Longitude,
	}
}

func PoolSettingsFromCommon(res *common.PoolSettings) (*PoolSettings, error) {
	if res == nil {
		return nil, nil //nolint: nilnil
	}

	shape, err := PoolShapeFromCommon(res.PoolShape)
	if err != nil {
		return nil, err
	}

	locationType, err := LocationTypeFromCommon(res.LocationType)
	if err != nil {
		return nil, err
	}

	usageType, err := UsageTypeFromCommon(res.UsageType)
	if err != nil {
		return nil, err
	}

	t, err := PoolTypeFromCommon(res.Type)
	if err != nil {
		return nil, err
	}

	return &PoolSettings{
		Type:         t,
		UsageType:    usageType,
		LocationType: locationType,
		Shape:        shape,
		Coordinates:  CoordinatesFromCommon(res.Coordinates),
	}, nil
}

func PoolTypeFromCommon(poolType common.PoolType) (PoolType, error) {
	switch poolType {
	case common.PoolTypeInfinity:
		return PoolTypeInfinity, nil
	case common.PoolTypeOverflow:
		return PoolTypeOverflow, nil
	case common.PoolTypeSkimmer:
		return PoolTypeSkimmer, nil
	default:
		return "", ErrUnknownPoolType
	}
}

func UsageTypeFromCommon(usageType common.UsageType) (UsageType, error) {
	switch usageType {
	case common.UsageTypePrivate:
		return UsageTypePrivate, nil
	case common.UsageTypePublic:
		return UsageTypeCommunity, nil
	case common.UsageTypeHoliday:
		return UsageTypeHoliday, nil
	default:
		return "", ErrUnknownUsageType
	}
}

func LocationTypeFromCommon(locationType common.LocationType) (LocationType, error) {
	switch locationType {
	case common.LocationTypeIndoor:
		return LocationTypeIndoor, nil
	case common.LocationTypeOutdoor:
		return LocationTypeOutdoor, nil
	default:
		return "", ErrUnknownLocationType
	}
}

func PoolShapeFromCommon(shape common.PoolShape) (PoolShape, error) {
	switch shape {
	case common.PoolShapeRectangle:
		return PoolShapeRectangle, nil
	case common.PoolShapeOval:
		return PoolShapeOval, nil
	case common.PoolShapeCircle:
		return PoolShapeCircle, nil
	case common.PoolShapeKidney:
		return PoolShapeKidney, nil
	case common.PoolShapeLShaped:
		return PoolShapeL, nil
	case common.PoolShapeTShaped:
		return PoolShapeT, nil
	case common.PoolShapeFreeForm:
		return PoolShapeFreeForm, nil
	default:
		return "", ErrUnknownPoolShape
	}
}

func CoordinatesFromCommon(coordinates *common.CoordinatesInput) *Coordinates {
	if coordinates == nil {
		return nil
	}

	return &Coordinates{
		Latitude:  coordinates.Latitude,
		Longitude: coordinates.Longitude,
	}
}
