package config

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"goblog/pkg/logger"
)

// Viper Viper 实例
var Viper *viper.Viper

// StrMap 简写 -- map[string]interface{}，用于方便书写
type StrMap map[string]interface{}

// init 初始化配置信息
func init() {
	// Initialize Viper
	Viper = viper.New()

	// 环境变量配置文件的的查找路径，相对于 main.go
	Viper.AddConfigPath(".")

	// 设置文件名
	Viper.SetConfigName(".env")

	// 设置配置文件类型，支持 "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	Viper.SetConfigType("env")

	// 开始读根目录下的 .env 文件，读不到会报错
	err := Viper.ReadInConfig()
	logger.LogError(err)

	// 设置环境变量前缀，用以区分 Go 的系统环境变量
	Viper.SetEnvPrefix("appenv")

	// Viper.AutomaticEnv() 会自动读取环境变量的值，用以覆盖掉配置文件中的值
	Viper.AutomaticEnv()
}

// Env 获取环境变量，支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return Get(envName, defaultValue[0])
	}
	return Get(envName)
}

// Add 添加配置
func Add(name string, configuration map[string]interface{}) {
	Viper.Set(name, configuration)
}

// Get 获取配置信息，支持通过 . 号来获取多级配置信息，例如：app.name
func Get(path string, defaultValue ...interface{}) interface{} {
	if !Viper.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return Viper.Get(path)
}

// GetString 获取字符串类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(Get(path, defaultValue...))
}

// GetInt 获取整型类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(Get(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(Get(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(Get(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(Get(path, defaultValue...))
}
