package types

import "time"

// Default settings
const (
	DefaultTimeout      = 15 * time.Second
	DefaultWriteTimeout = 5 * time.Second
	MaxBufferSize       = 2 << 24
)
