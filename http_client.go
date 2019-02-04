package client

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"github.com/nilorg/sdk/log"

	"github.com/robfig/cron"
)

const cargoboatConfigVersionKey = "cargoboat.config.version"

// CargoboatClient 客户端
type CargoboatClient struct {
	httpClient         *http.Client
	cron               *cron.Cron
	log                log.Logger
	lock               sync.RWMutex
	config             map[string]interface{}
	baseURL            *url.URL
	configVersion      int64
	username, password string
	cronSpec           string
}

// NewCargoboatClient 创建 redis客户端
func NewCargoboatClient(log log.Logger, baseURL, username, password, cronSpec string) Clienter {
	client := &CargoboatClient{
		httpClient: &http.Client{},
		cron:       cron.New(),
		log:        log,
		lock:       sync.RWMutex{},
		config:     make(map[string]interface{}),
		username:   username,
		password:   password,
		cronSpec:   "*/5 * * * * ?",
	}
	if cronSpec != "" {
		client.cronSpec = cronSpec
	}

	u, _ := url.Parse(baseURL)
	u.Path = path.Join(u.Path, "client")
	client.baseURL = u
	client.init()
	return client
}

// urlJoin url 拼接
func (c *CargoboatClient) urlJoin(uri string) string {
	u := *c.baseURL
	u.Path = path.Join(u.Path, uri)
	return u.String()
}

// configItem ...
type configItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// configResult ...
type configResult struct {
	Version int64        `json:"version"`
	Configs []configItem `json:"configs"`
}

// init 初始化
func (c *CargoboatClient) init() {
	var req *http.Request
	var err error
	req, err = http.NewRequest(http.MethodGet, c.urlJoin("/configs"), nil)
	if err != nil {
		c.log.Errorln(err)
		return
	}
	var resp *http.Response
	resp, err = c.do(req)
	if err != nil {
		c.log.Errorln(err)
		return
	}
	if resp.StatusCode != 200 {
		c.log.Debugf("check version resp status code:%d", resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	jsonDecode := json.NewDecoder(resp.Body)
	configResult := configResult{}
	err = jsonDecode.Decode(&configResult)
	if err != nil {
		c.log.Errorln(err)
		return
	}
	// 设置配置项
	c.log.Debugf("init config,local version:%d/server version:%d", c.configVersion, configResult.Version)
	if c.configVersion != configResult.Version {
		c.set(configResult.Configs...)
		c.configVersion = configResult.Version
	}
}

// do 执行HTTP请求
func (c *CargoboatClient) do(req *http.Request) (response *http.Response, err error) {
	req.SetBasicAuth(c.username, c.password)
	response, err = c.httpClient.Do(req)
	return
}

// set 批量设置配置项
func (c *CargoboatClient) set(value ...configItem) {
	defer c.lock.Unlock()
	c.lock.Lock()

	for _, v := range value {
		if v.Key == cargoboatConfigVersionKey {
			continue
		}
		c.config[v.Key] = v.Value
	}
}

// getConfig 获取配置
func (c *CargoboatClient) getConfig(key string) interface{} {
	defer c.lock.RUnlock()
	c.lock.RLock()
	return c.config[key]
}

// versionResult 版本结果
type versionResult struct {
	Version int64 `json:"version"`
}

// checkVersion 检查版本
func (c *CargoboatClient) checkVersion() {
	var req *http.Request
	var err error
	req, err = http.NewRequest(http.MethodGet, c.urlJoin("/version"), nil)
	if err != nil {
		c.log.Errorln(err)
		return
	}
	var resp *http.Response
	resp, err = c.do(req)
	if err != nil {
		c.log.Errorln(err)
		return
	}

	if resp.StatusCode != 200 {
		c.log.Debugf("check version resp status code:%d", resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	jsonDecode := json.NewDecoder(resp.Body)
	result := versionResult{}
	err = jsonDecode.Decode(&result)
	if err != nil {
		c.log.Errorln(err)
		return
	}
	// 设置配置项
	c.log.Debugf("check version,local version:%d/server version:%d", c.configVersion, result.Version)
	if c.configVersion != result.Version {
		c.init()
	}
}

// WatchConfig 监听配置
func (c *CargoboatClient) WatchConfig() {
	err := c.cron.AddFunc(c.cronSpec, c.checkVersion)
	if err != nil {
		c.log.Errorf("WatchConfig AddFunc:%v", err)
	}
	c.cron.Start()
}

// Get return value as a interface{}.
func (c *CargoboatClient) Get(key string) interface{} {
	value := c.getConfig(key)
	c.log.Debugf("Get %s Value:%v", key, value)
	return value
}

// GetBool return value as a bool.
func (c *CargoboatClient) GetBool(key string) bool {
	value := c.getConfig(key)
	c.log.Debugf("GetBool %s Value:%v", key, value)
	return value.(bool)
}

// GetFloat64 return value as a float64.
func (c *CargoboatClient) GetFloat64(key string) float64 {
	value := c.getConfig(key)
	c.log.Debugf("GetFloat64 %s Value:%v", key, value)
	return value.(float64)
}

// GetInt return value as a int.
func (c *CargoboatClient) GetInt(key string) int {
	value := c.getConfig(key)
	c.log.Debugf("GetInt %s Value:%v", key, value)
	return value.(int)
}

// GetString return value as a string.
func (c *CargoboatClient) GetString(key string) string {
	value := c.getConfig(key)
	c.log.Debugf("GetString %s Value:%v", key, value)
	return value.(string)
}

// GetTime return value as a time.Time.
func (c *CargoboatClient) GetTime(key string) time.Time {
	value := c.getConfig(key)
	c.log.Debugf("GetTime %s Value:%v", key, value)
	return value.(time.Time)
}

// GetDuration return value as a time.Duration.
func (c *CargoboatClient) GetDuration(key string) time.Duration {
	value := c.getConfig(key)
	c.log.Debugf("GetDuration %s Value:%v", key, value)
	return value.(time.Duration)
}

// IsExist return key is exist.
func (c *CargoboatClient) IsExist(key string) bool {
	defer c.lock.RUnlock()
	c.lock.RLock()
	_, ok := c.config[key]
	return ok
}

// Close 关闭
func (c *CargoboatClient) Close() error {
	c.cron.Stop()
	return nil
}
