package client

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

// Config ...
type Config map[string]interface{}

// ConfigItem ...
type ConfigItem struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// RedisCargoboatClient redis 客户端
type RedisCargoboatClient struct {
	redisClient       *redis.Client
	log               Logger
	appKey, appSecret string
	lock              sync.RWMutex
	config            Config
}

// NewRedisCargoboatClient 创建 redis客户端
func NewRedisCargoboatClient(appKey, appSecret string, log Logger, options *redis.Options) Clienter {
	redisClient := redis.NewClient(options)
	pong, err := redisClient.Ping().Result()
	if err != nil {
		log.Errorf("redis client ping:%v", err)
		return nil
	}
	log.Infof("redis client pong:%s", pong)

	client := &RedisCargoboatClient{
		log:         log,
		redisClient: redisClient,
		appKey:      appKey,
		appSecret:   appSecret,
		lock:        sync.RWMutex{},
		config:      Config{},
	}
	client.init()
	return client
}

// init 初始化
func (c *RedisCargoboatClient) init() {
	c.log.Infoln("cargoboat init ...")
	configKey := fmt.Sprintf("cargoboat_%s_%s_config_keys", c.appKey, c.appSecret)
	// 获取 set集合数量
	icmd := c.redisClient.SCard(configKey)
	if icmd.Val() == 0 {
		c.log.Warningln("cargoboat server config keys count:0")
		return
	}
	// 查找set集合中的config key array
	scmd := c.redisClient.SScan(configKey, 0, "", icmd.Val())
	resultArray, _ := scmd.Val()
	items := make([]*ConfigItem, 0)
	for _, key := range resultArray {
		scmd := c.redisClient.Get(c.wrapKey(key))
		if scmd.Err() == redis.Nil {
			continue
		}
		items = append(items, &ConfigItem{
			Key:   key,
			Value: scmd.Val(),
		})
	}
	// 设置配置项
	c.setConfig(items...)
}

// wrapKey 包装Key
func (c *RedisCargoboatClient) wrapKey(key string) string {
	return fmt.Sprintf("cargoboat_%s_%s_%s", c.appKey, c.appSecret, key)
}

func (c *RedisCargoboatClient) setConfig(items ...*ConfigItem) {
	defer c.lock.Unlock()
	c.lock.Lock()
	for _, item := range items {
		c.log.Debugf("set config %s:%v", item.Key, item.Value)
		c.config[item.Key] = item.Value
	}
}

// String2ConfigItem string to ConfigItem.
func String2ConfigItem(msg string) *ConfigItem {
	if msg == "" {
		return nil
	}
	citem := &ConfigItem{}
	if err := json.Unmarshal([]byte(msg), citem); err != nil {
		return nil
	}
	return citem
}

func (c *RedisCargoboatClient) subConfig(msg string) {
	c.log.Debugf("cargoboat server pub msg:%s", msg)
	item := String2ConfigItem(msg)
	if item != nil {
		c.setConfig(item)
	}
}

func (c *RedisCargoboatClient) getConfig(key string) interface{} {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.config[key]; ok {
		return value
	}
	return nil
}

// WatchConfig 监听配置
func (c *RedisCargoboatClient) WatchConfig() {
	subKey := fmt.Sprintf("cargoboat_%s_%s_channel", c.appKey, c.appSecret)
	pubsub := c.redisClient.Subscribe(subKey)
	// 订阅内容之前，等待确认订阅已创建。
	receiveResult, err := pubsub.Receive()
	if err != nil {
		c.log.Errorf("cargoboat client subscribe error:%v", err)
		return
	}
	c.log.Infof("cargoboat client subscribe receive:%v", receiveResult)
	defer pubsub.Close()
	// 接收消息的通道。
	ch := pubsub.Channel()
	for {
		msg, ok := <-ch
		if !ok {
			continue
		}
		go c.subConfig(msg.Payload)
	}
}

// Get return value as a interface{}.
func (c *RedisCargoboatClient) Get(key string) interface{} {
	value := c.getConfig(key)
	c.log.Debugf("Get %s Value:%v", key, value)
	return value
}

// GetBool return value as a bool.
func (c *RedisCargoboatClient) GetBool(key string) bool {
	value := c.getConfig(key)
	c.log.Debugf("GetBool %s Value:%v", key, value)
	return value.(bool)
}

// GetFloat64 return value as a float64.
func (c *RedisCargoboatClient) GetFloat64(key string) float64 {
	value := c.getConfig(key)
	c.log.Debugf("GetFloat64 %s Value:%v", key, value)
	return value.(float64)
}

// GetInt return value as a int.
func (c *RedisCargoboatClient) GetInt(key string) int {
	value := c.getConfig(key)
	c.log.Debugf("GetInt %s Value:%v", key, value)
	return value.(int)
}

// GetString return value as a string.
func (c *RedisCargoboatClient) GetString(key string) string {
	value := c.getConfig(key)
	c.log.Debugf("GetString %s Value:%v", key, value)
	return value.(string)
}

// GetStringMap return value as a map[string]interface{}.
func (c *RedisCargoboatClient) GetStringMap(key string) map[string]interface{} {
	value := c.getConfig(key)
	c.log.Debugf("GetStringMap %s Value:%v", key, value)
	return value.(map[string]interface{})
}

// GetStringMapString return value as a map[string]string.
func (c *RedisCargoboatClient) GetStringMapString(key string) map[string]string {
	value := c.getConfig(key)
	c.log.Debugf("GetStringMapString %s Value:%v", key, value)
	return value.(map[string]string)
}

// GetStringSlice return value as a []string.
func (c *RedisCargoboatClient) GetStringSlice(key string) []string {
	value := c.getConfig(key)
	c.log.Debugf("GetStringSlice %s Value:%v", key, value)
	return value.([]string)
}

// GetTime return value as a time.Time.
func (c *RedisCargoboatClient) GetTime(key string) time.Time {
	value := c.getConfig(key)
	c.log.Debugf("GetTime %s Value:%v", key, value)
	return value.(time.Time)
}

// GetDuration return value as a time.Duration.
func (c *RedisCargoboatClient) GetDuration(key string) time.Duration {
	value := c.getConfig(key)
	c.log.Debugf("GetDuration %s Value:%v", key, value)
	return value.(time.Duration)
}

// IsExist return key is exist.
func (c *RedisCargoboatClient) IsExist(key string) bool {
	defer c.lock.RUnlock()
	c.lock.RLock()
	_, ok := c.config[key]
	return ok
}

// AllSettings return value as a map[string]interface{}.
func (c *RedisCargoboatClient) AllSettings() map[string]interface{} {
	defer c.lock.RUnlock()
	c.lock.RLock()
	addr := &c.config
	return *addr
}
