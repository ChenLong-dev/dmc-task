package core

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var Cfg *Config

const (
	DefaultCfgPath = "./conf/conf.yaml"
)

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	ApiServer  ApiServerConfig  `mapstructure:"apiserver"`
	GrpcServer GrpcServerConfig `mapstructure:"grpcserver"`
	Logx       LogxConfig       `mapstructure:"logx"`
	MySQL      MySqlConfig      `mapstructure:"mysql"`
}

type AppConfig struct {
	Name          string `mapstructure:"name"`
	Mode          string `mapstructure:"mode"`
	Version       string `mapstructure:"version"`
	IsDistributed bool   `mapstructure:"is_distributed"`
}

type ApiServerConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
}

type GrpcServerConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
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

func ConfigInit(path string, cfg interface{}) error {
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
	err = viper.Unmarshal(cfg)
	if err != nil {
		fmt.Printf("viper unmarshal is failed! err:%v\n", err)
		return err
	}
	fmt.Printf("config init success! cfg:%+v\n", cfg)

	return nil
}
