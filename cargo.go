package cargo

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cargoboat/cargo/client"
)

var cargoboat client.Clienter

var (
	// AppKey ...
	AppKey string
	// AppSecret ...
	AppSecret string
	// ServerAddr ...
	ServerAddr string
)

// Init 初始化
func Init() {
	cargoboat = client.NewCargoboatClient(logrus.New(), ServerAddr, AppKey, AppSecret, "")
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

// GetTime return value as a time.Time.
func GetTime(key, timeLayout string) (time.Time, error) {
	return cargoboat.GetTime(key, timeLayout)
}

// GetDuration return value as a time.Duration.
func GetDuration(key string) time.Duration {
	return cargoboat.GetDuration(key)
}

// GetEnv return value as a interface{}.
func GetEnv(key string) interface{} {
	return cargoboat.GetEnv(key)
}

// IsExist return key is exist.
func IsExist(key string) bool {
	return cargoboat.IsExist(key)
}
