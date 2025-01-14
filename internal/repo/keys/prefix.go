package keys

type Prefix [2]byte

var (
	// Prefixes for keys in the new format
	versionPrefix      Prefix = [2]byte{0x00, 0x00}
	poolPrefix         Prefix = [2]byte{0x00, 0x01}
	userPrefix         Prefix = [2]byte{0x00, 0x02}
	devicePrefix       Prefix = [2]byte{0x00, 0x03}
	notificationPrefix Prefix = [2]byte{0x00, 0x04}
	poolSettingsPrefix Prefix = [2]byte{0x00, 0x05}

	// Prefixes for keys of indexes
	userIndexAppleID       Prefix = [2]byte{0x01, 0x00}
	userIndexNotifications Prefix = [2]byte{0x01, 0x01}
	userIndexDevices       Prefix = [2]byte{0x01, 0x02}
	userIndexPool          Prefix = [2]byte{0x01, 0x03}
)
