package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

type Config struct {
	Logger   LoggerConf `mapstructure:"logger"`
	App      AppConf    `mapstructure:"app"`
	Database DBConf     `mapstructure:"database"`
}

type LoggerConf struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type AppConf struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	Storage string `mapstructure:"storage"`
}

type DBConf struct {
	Host   string `mapstructure:"host"`
	Port   string `mapstructure:"port"`
	DBName string `mapstructure:"dbname"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
}

func NewConfig() Config {
	var config Config

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("init config error: %v", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Sprintf("init config error: %v", err))
	}

	return config
}
