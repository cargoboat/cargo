package client

import (
	"time"
)

// Clienter ...
type Clienter interface {
	WatchConfig()
	Get(key string) interface{}
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetString(key string) string
	GetTime(key, timeLayout string) (time.Time, error)
	GetDuration(key string) time.Duration
	GetEnv(key string) interface{}
	IsExist(key string) bool
	Close() error
}
