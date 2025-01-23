package core

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var Cfg *Config

const (
	DefaultCfgPath = "./conf/conf.yaml"
)

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Server ServerConfig `mapstructure:"server"`
	Logx   LogxConfig   `mapstructure:"logx"`
	MySQL  MySqlConfig  `mapstructure:"mysql"`
}

type AppConfig struct {
	Name          string `mapstructure:"name"`
	Mode          string `mapstructure:"mode"`
	Version       string `mapstructure:"version"`
	IsDistributed bool   `mapstructure:"is_distributed"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LogxConfig struct {
	Mode       string `mapstructure:"mode"`
	Encoding   string `mapstructure:"encoding"`
	Path       string `mapstructure:"path"`
	Level      string `mapstructure:"level"`
	KeepDays   int    `mapstructure:"keep_days"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxSize    int    `mapstructure:"max_size"`
	Rotation   string `mapstructure:"rotation"`
}

type MySqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
	Timeout  int    `mapstructure:"Timeout"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

func ConfigInit(path string) error {
	if !IsExist(path) {
		path = DefaultCfgPath
	}
	fmt.Printf("config init path:%s\n", path)
	basePath := filepath.Base(path)
	dirPath := filepath.Dir(path)
	baseList := strings.Split(basePath, ".") // ["conf", "yaml"]
	if len(baseList) != 2 {
		return fmt.Errorf("file's format is error! path:%s", path)
	}
	viper.SetConfigName(baseList[0])
	viper.SetConfigType(baseList[1])
	viper.AddConfigPath(dirPath)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper read is failed! err:%v\n", err)
		return err
	}
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Printf("viper unmarshal is failed! err:%v\n", err)
		return err
	}
	fmt.Printf("config init success! cfg:%+v\n", cfg)
	Cfg = &cfg
	return nil
}
