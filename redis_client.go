package client

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/nilorg/sdk/convert"
)

// RedisCargoboatClient redis 客户端
type RedisCargoboatClient struct {
	redisClient       *redis.Client
	log               Logger
	appKey, appSecret string
}

// NewRedisCargoboatClient 创建 redis客户端
func NewRedisCargoboatClient(appKey, appSecret string, log Logger, options *redis.Options) Clienter {
	client := redis.NewClient(options)
	pong, err := client.Ping().Result()
	if err != nil {
		log.Errorf("redis client ping:%v", err)
		return nil
	}
	log.Infof("redis client pong:%s", pong)
	return &RedisCargoboatClient{
		log:         log,
		redisClient: client,
		appKey:      appKey,
		appSecret:   appSecret,
	}
}

// wrapKey 包装Key
func (c *RedisCargoboatClient) wrapKey(key string) string {
	return fmt.Sprintf("%s_%s_%s", c.appKey, c.appSecret, key)
}

// WatchConfig 监听配置
func (c *RedisCargoboatClient) WatchConfig() {

}

// Get return value as a interface{}.
func (c *RedisCargoboatClient) Get(key string) interface{} {
	value, err := c.redisClient.Get(c.wrapKey(key)).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		c.log.Errorf("Get %s Error:%v", key, err)
		return nil
	}
	c.log.Debugf("Get %s Value:%s", value)
	return value
}

// GetBool return value as a bool.
func (c *RedisCargoboatClient) GetBool(key string) bool {
	value, err := c.redisClient.Get(c.wrapKey(key)).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		c.log.Errorf("GetBool %s Error:%v", key, err)
		return false
	}
	if value == "1" || value == "true" {
		return true
	}
	c.log.Debugf("GetBool %s Value:%s", value)
	return false
}

// GetFloat64 return value as a float64.
func (c *RedisCargoboatClient) GetFloat64(key string) float64 {
	value, err := c.redisClient.Get(c.wrapKey(key)).Result()
	if err == redis.Nil {
		return 0
	} else if err != nil {
		c.log.Errorf("GetFloat64 %s Error:%v", key, err)
		return 0
	}
	c.log.Debugf("GetFloat64 %s Value:%s", value)
	return convert.ToFloat64(value)
}

// GetInt return value as a int.
func (c *RedisCargoboatClient) GetInt(key string) int {
	value, err := c.redisClient.Get(c.wrapKey(key)).Result()
	if err == redis.Nil {
		return 0
	} else if err != nil {
		c.log.Errorf("GetInt %s Error:%v", key, err)
		return 0
	}
	c.log.Debugf("GetInt %s Value:%s", value)
	return convert.ToInt(value)
}

// GetString return value as a string.
func (c *RedisCargoboatClient) GetString(key string) string {
	value, err := c.redisClient.Get(c.wrapKey(key)).Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		c.log.Errorf("GetString %s Error:%v", key, err)
		return ""
	}
	c.log.Debugf("GetString %s Value:%s", value)
	return value
}

// GetStringMap return value as a map[string]interface{}.
func (c *RedisCargoboatClient) GetStringMap(key string) map[string]interface{} {
	return nil
}

// GetStringMapString return value as a map[string]string.
func (c *RedisCargoboatClient) GetStringMapString(key string) map[string]string {
	return nil
}

// GetStringSlice return value as a []string.
func (c *RedisCargoboatClient) GetStringSlice(key string) []string {
	return nil
}

// GetTime return value as a time.Time.
func (c *RedisCargoboatClient) GetTime(key string) time.Time {

	return time.Now()
}

// GetDuration return value as a time.Duration.
func (c *RedisCargoboatClient) GetDuration(key string) time.Duration {
	return 0
}

// IsExist return key is exist.
func (c *RedisCargoboatClient) IsExist(key string) bool {

	return false
}

// AllSettings return value as a map[string]interface{}.
func (c *RedisCargoboatClient) AllSettings() map[string]interface{} {
	return nil
}
