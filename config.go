package goconf

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Config 结构声明
type Config struct {
	Key         string
	Default     string
	Required    bool
	Description string
}

// New New
func New() *ConfigLoader {
	configService := new(ConfigLoader)
	configService.data = make(map[string]*string)
	configService.standards = make(map[string]Config)
	return configService
}

// ConfigLoader 配置
type ConfigLoader struct {
	loaded  bool // 加载完成
	checked bool // 检查完成

	data      map[string]*string
	standards map[string]Config
	sync.RWMutex
}

// Set 设置一个配置
func (c *ConfigLoader) Set(key string, value string) {
	c.RLock()
	defer c.RUnlock()
	c.data[key] = &value
}

// Get 获取一个配置
func (c *ConfigLoader) Get(key string) (value string, err error) {
	if !c.loaded {
		if err = c.load(); err != nil {
			return "", err
		}
	}

	if !c.checked {
		if err = c.check(); err != nil {
			return "", err
		}
	}

	if _, ok := c.standards[key]; !ok {
		return "", fmt.Errorf("config: %s is not declared", key)
	}

	if c.data[key] == nil {
		return "", fmt.Errorf("config: %s is nil", key)
	}

	return *c.data[key], nil
}

func (c *ConfigLoader) GetInt(key string) (value int, err error) {
	stringValue, err := c.Get(key)
	if err != nil {
		return 0, err
	}

	value, err = strconv.Atoi(stringValue)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// MustGet 获取一个配置
func (c *ConfigLoader) MustGet(key string) (value string) {
	value, err := c.Get(key)
	if err != nil {
		panic(err)
	}

	return value
}

// MustGet 获取一个配置
func (c *ConfigLoader) MustGetInt(key string) (value int) {
	value, err := c.GetInt(key)
	if err != nil {
		panic(err)
	}

	return value
}

// Declare 添加配置定义
func (c *ConfigLoader) Declare(key string, deft string, required bool, description string) {
	c.RLock()
	defer c.RUnlock()

	// 恢复为未 check 状态
	c.checked = false

	// 检查 key 格式
	c.keyCheck(key)

	// 记录
	stan := new(Config)
	stan.Key = key
	stan.Default = deft
	stan.Required = required
	stan.Description = description
	c.standards[stan.Key] = *stan

	// 注册 flag
	var value string
	value = stan.Default
	c.data[stan.Key] = &value
	flag.StringVar(&value, key, stan.Default, description)
}

// Load 加载数据
func (c *ConfigLoader) Load() (err error) {
	if !c.loaded {
		if err := c.load(); err != nil {
			return err
		}
	}

	if !c.checked {
		if err := c.check(); err != nil {
			return err
		}
	}

	return nil
}

// MustLoad 确保加载完成
func (c *ConfigLoader) MustLoad() {
	if err := c.Load(); err != nil {
		panic(err)
	}
}

// check 检查加载到的数据
func (c *ConfigLoader) load() (err error) {
	c.loadEnv()  // 先加载环境变量
	c.loadFlag() // 再加载命令行参数
	c.loaded = true
	return nil
}

// check 检查加载到的数据
func (c *ConfigLoader) check() (err error) {
	for _, standard := range c.standards {
		if standard.Required && c.data[standard.Key] == nil {
			return fmt.Errorf("config: %s is required, %s", standard.Key, standard.Description)
		}
	}

	c.checked = true
	return nil
}

// loadFlag 加载启动命令参数
func (c *ConfigLoader) loadFlag() {
	c.RLock()
	defer c.RUnlock()
	c.checked = false
	cache := make(map[string]*string)

	flag.Parse()

	for key, value := range cache {
		if value != nil && *value != "" {
			c.data[key] = value
		}
	}
}

// 加载环境变量
func (c *ConfigLoader) loadEnv() {
	c.RLock()
	defer c.RUnlock()
	c.checked = false
	for key, standard := range c.standards {
		value := os.Getenv(standard.Key)
		if value != "" {
			c.data[key] = &value
		}

		valueByUpperKey := os.Getenv(strings.ToUpper(standard.Key))
		if valueByUpperKey != "" {
			c.data[key] = &valueByUpperKey
		}
	}
}

// MustKeyCheck 检查
func (c *ConfigLoader) keyCheck(key string) {
	if key == "" {
		panic(fmt.Errorf("config: 配置 key 名不允许为空"))
	}

	matched, err := regexp.MatchString("^[a-zA-Z_]*$", key)

	if err != nil {
		panic(fmt.Errorf("config: 配置 key 名检查错误: %v", err))
	}

	if !matched {
		panic(fmt.Errorf("config: 配置 key 名仅允许大小写字母、下划线"))
	}
}
