package cc

const (
	UNINITIALIZED_PANIC_STATE = "[panic state] setup/un-initialized"
	INVALID_ENV_PANIC_STATE   = "[panic state] invalid environment variable"
	INVALID_DB_CONFIG         = "invalid DB config"
	INVALID_HOST_CONFIG       = "invalid host config"

	// Redis Error
	BUSY_GROUP_ALREADY_EXISTS = "BUSYGROUP Consumer Group name already exists"
	// Handler Errors
	NO_PAYLOAD_FOUND = "no payload found"
)
