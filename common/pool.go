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

type PoolMetadata struct {
	CleanerUserID uuid.UUID
}
