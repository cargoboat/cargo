package client

import (
	"time"
)

// Clienter ...
type Clienter interface {
	WatchConfig()
	Get(key string) interface{}
	GetBool(key string) bool
	GetFloat32(key string) float32
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetUint(key string) uint
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	GetString(key string) string
	GetTime(key, timeLayout string) (time.Time, error)
	GetDuration(key string) time.Duration
	GetEnv(key string) interface{}
	IsExist(key string) bool
	Close() error
}
