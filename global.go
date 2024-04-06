package goconf

// globalService Config
var globalService *ConfigLoader

func init() {
	globalService = New()
}

func Set(key string, value string) {
	globalService.Set(key, value)
}

// Get 获取一个配置
func Get(key string) (value string, err error) {
	return globalService.Get(key)
}

func GetInt(key string) (value int, err error) {
	return globalService.GetInt(key)
}

// MustGet 获取一个配置
func MustGet(key string) (value string) {
	return globalService.MustGet(key)
}

func MustGetInt(key string) (value int) {
	return globalService.MustGetInt(key)
}

// Declare 设置定义
func Declare(key string, deft string, required bool, description string) {
	globalService.Declare(key, deft, required, description)
}

func Load() (error) {
	return globalService.Load()
}

func MustLoad() {
	globalService.MustLoad()
}
