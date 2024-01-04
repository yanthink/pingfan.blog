package config

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {
	// 默认配置文件
	viper.SetConfigFile("config.toml")
	// 配置类型，支持 "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config.toml 配置文件读取失败: %s", err))
	}

	if _, err := os.Stat("config.local.toml"); err == nil {
		// 加载私有配置文件
		viper.SetConfigName("config.local.toml")
		// 环境变量配置文件查找的路径，相对于 main.go
		viper.AddConfigPath(".")
		if err := viper.MergeInConfig(); err != nil {
			panic(fmt.Errorf("合并私有配置失败: %s", err))
		}
	}

	viper.SetEnvPrefix("BLOG")
	// 将所有的 . 和 - 替换成 _，GetString("mysql.host") 会优先读取 BLOG_MYSQL_HOST 环境变量的值
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// 加载应用配置
	loadAppConfig()
	// 加载数据库配置
	loadDatabaseConfig()
	// 加载jwt配置
	loadJwtConfig()
	// 加载redis配置
	loadRedisConfig()
	// 加载cache配置
	loadCacheConfig()
	// 加载snowflake配置
	loadSnowflakeConfig()
	// 加载storage配置
	loadStorageConfig()
	// 加载mail配置
	loadMailConfig()
	// 加载zincsearch配置
	loadZincsearchConfig()
	// 加载微信小程序配置
	loadMiniProgramConfig()
}

// Get 获取配置项，允许使用点式获取，如：app.name
func Get(path string, defaultValue ...any) any {
	// 不存在的情况
	if !viper.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	return viper.Get(path)
}

// GetString 获取 String 类型的配置信息
func GetString(path string, defaultValue ...any) string {
	return cast.ToString(Get(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...any) int {
	return cast.ToInt(Get(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...any) int64 {
	return cast.ToInt64(Get(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...any) uint {
	return cast.ToUint(Get(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...any) bool {
	return cast.ToBool(Get(path, defaultValue...))
}
