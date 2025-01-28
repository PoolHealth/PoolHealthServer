package keys

import (
	"encoding/hex"

	"github.com/google/uuid"
)

type Builder interface {
	Version() string
	Pools() string
	UserPools(userID uuid.UUID) string
	UsersPools() string
	Pool(id uuid.UUID) string
	Users() string
	User(id uuid.UUID) string
	UserByAppleID(id string) string
	Device(id uuid.UUID) string
	DevicesByUserID(id uuid.UUID) string
	Notification(id uuid.UUID) string
	NotificationByUserID(id uuid.UUID) string
	PoolSettings(poolID uuid.UUID) string
}

type builder struct {
}

func (b builder) PoolSettings(poolID uuid.UUID) string {
	return appendUUID(poolSettingsPrefix, poolID)
}

func (b builder) UserPools(userID uuid.UUID) string {
	return appendUUID(userIndexPool, userID)
}

func (b builder) UsersPools() string {
	return hex.EncodeToString(userIndexPool[:])
}

func (b builder) Users() string {
	return hex.EncodeToString(userPrefix[:])
}

func (b builder) User(id uuid.UUID) string {
	return appendUUID(userPrefix, id)
}

func (b builder) UserByAppleID(id string) string {
	return hex.EncodeToString(userIndexAppleID[:]) + id
}

func (b builder) Device(id uuid.UUID) string {
	return appendUUID(devicePrefix, id)
}

func (b builder) DevicesByUserID(id uuid.UUID) string {
	return appendUUID(userIndexDevices, id)
}

func (b builder) Notification(id uuid.UUID) string {
	return appendUUID(notificationPrefix, id)
}

func (b builder) NotificationByUserID(id uuid.UUID) string {
	return appendUUID(userIndexNotifications, id)
}

func (b builder) Pool(id uuid.UUID) string {
	return appendUUID(poolPrefix, id)
}

func (b builder) Pools() string {
	return hex.EncodeToString(poolPrefix[:])
}

func (b builder) Version() string {
	return hex.EncodeToString(versionPrefix[:])
}

func NewBuilder() Builder {
	return &builder{}
}

func appendPrefix(prefix [2]byte, data string) string {
	return hex.EncodeToString(prefix[:]) + data
}

func appendUUID(prefix [2]byte, id uuid.UUID) string {
	return appendPrefix(prefix, id.String())
}
