package client

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

// Clienter ...
type Clienter interface {
	WatchConfig()
	Get(key string) interface{}
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	IsExist(key string) bool
	AllSettings() map[string]interface{}
}

var cargoboat Clienter

var (
	// AppKey ...
	AppKey string
	// AppSecret ...
	AppSecret string
	// ClientAddr ...
	ClientAddr string
	// ClientPass ...
	ClientPass string
)

// Init 初始化
func Init() {
	cargoboat = NewRedisCargoboatClient(AppKey, AppSecret, logrus.New(), &redis.Options{
		Addr:     ClientAddr,
		Password: ClientPass,
		DB:       0,
	})
}

// WatchConfig 监听配置
func WatchConfig() {
	cargoboat.WatchConfig()
}

// Get return value as a interface{}.
func Get(key string) interface{} {
	return cargoboat.Get(key)
}

// GetBool return value as a bool.
func GetBool(key string) bool {
	return cargoboat.GetBool(key)
}

// GetFloat64 return value as a float64.
func GetFloat64(key string) float64 {
	return cargoboat.GetFloat64(key)
}

// GetInt return value as a int.
func GetInt(key string) int {
	return cargoboat.GetInt(key)
}

// GetString return value as a string.
func GetString(key string) string {
	return cargoboat.GetString(key)
}

// GetStringMap return value as a map[string]interface{}.
func GetStringMap(key string) map[string]interface{} {
	return cargoboat.GetStringMap(key)
}

// GetStringMapString return value as a map[string]string.
func GetStringMapString(key string) map[string]string {
	return cargoboat.GetStringMapString(key)
}

// GetStringSlice return value as a []string.
func GetStringSlice(key string) []string {
	return cargoboat.GetStringSlice(key)
}

// GetTime return value as a time.Time.
func GetTime(key string) time.Time {
	return cargoboat.GetTime(key)
}

// GetDuration return value as a time.Duration.
func GetDuration(key string) time.Duration {
	return cargoboat.GetDuration(key)
}

// IsExist return key is exist.
func IsExist(key string) bool {
	return cargoboat.IsExist(key)
}

// AllSettings return value as a map[string]interface{}.
func AllSettings() map[string]interface{} {
	return cargoboat.AllSettings()
}
