package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

var Conf config

type config struct {
	Logs   *LogConfig `mapstructure:"logs" json:"logs"`
	App    *AppConfig `mapstructure:"apps" json:"apps"`
	Jwt    *Jwt       `mapstructure:"jwt" json:"jwt"`
	System *System    `mapstructure:"system" json:"system"`
}

func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("读取应用目录失败:%s", err))
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/")

	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败:%s", err))
	}

	// 将配置映射到全局变量Conf
	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(fmt.Errorf("初始化配置文件失败:%s", err))
	}
}

type LogConfig struct {
	Level      zapcore.Level `mapstructure:"level" json:"level"`
	Path       string        `mapstructure:"path" json:"path"`
	MaxSize    int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	Maxage     int           `mapstructure:"max-age" json:"maxAge"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
}

type AppConfig struct {
	Mode string `mapstructure:"mode" json:"mode"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type Jwt struct {
	Realm      string `mapstructure:"realm" json:"realm"`
	Key        string `mapstructure:"key" json:"key"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int    `mapstructure:"max-refresh" json:"maxRefresh"`
}

type System struct {
	RSAPrivateBytes []byte `mapstructure:"_" json:"_"`
}
