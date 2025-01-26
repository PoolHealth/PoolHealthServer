package common

import "github.com/google/uuid"

type Migration struct {
	ID     uuid.UUID
	Status MigrationStatus
}

type MigrationStatus uint

const (
	// MigrationStatus enum
	MigrationStatusUnknown MigrationStatus = iota
	MigrationStatusPending
	MigrationsStatusCompleted
	MigrationStatusFailed
)
