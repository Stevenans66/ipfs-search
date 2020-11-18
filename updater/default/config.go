package updater

import (
	"time"
)

// Config is the configuration for a DefaultUpdater.
type Config struct {
	// MinAge is the minimum age a resource should have to have its LastSeen updated.
	MinAge time.Duration
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		MinAge: time.Hour,
	}
}
