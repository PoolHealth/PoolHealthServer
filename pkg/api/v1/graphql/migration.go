package graphql

import (
	rootCommon "github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/api/v1/common"
)

func MigrationFromCommon(m rootCommon.Migration) *Migration {
	return &Migration{
		ID:     common.ID(m.ID),
		Status: MigrationStatusFromCommon(m.Status),
	}
}

func MigrationStatusFromCommon(s rootCommon.MigrationStatus) MigrationStatus {
	switch s {
	case rootCommon.MigrationStatusPending:
		return MigrationStatusPending
	case rootCommon.MigrationStatusUnknown:
		return MigrationStatusUnknown
	case rootCommon.MigrationsStatusCompleted:
		return MigrationStatusDone
	case rootCommon.MigrationStatusFailed:
		return MigrationStatusFailed
	}

	return MigrationStatusUnknown
}
