package core

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var FrontendCfg *FrontendConfig

const (
	DefaultFrontendCfgPath = "./conf/conf.yaml"
)

type FrontendConfig struct {
	App       FrontendAppConfig       `mapstructure:"app"`
	ApiServer FrontendApiServerConfig `mapstructure:"apiserver"`
	Logx      FrontendLogxConfig      `mapstructure:"logx"`
}

type FrontendAppConfig struct {
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
}

type FrontendApiServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type FrontendLogxConfig struct {
	Mode       string `mapstructure:"mode"`
	Encoding   string `mapstructure:"encoding"`
	Path       string `mapstructure:"path"`
	Level      string `mapstructure:"level"`
	KeepDays   int    `mapstructure:"keep_days"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxSize    int    `mapstructure:"max_size"`
	Rotation   string `mapstructure:"rotation"`
}

func FrontendConfigInit(path string) error {
	if !IsExist(path) {
		path = DefaultFrontendCfgPath
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
	var cfg FrontendConfig
	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Printf("viper unmarshal is failed! err:%v\n", err)
		return err
	}
	fmt.Printf("config init success! cfg:%+v\n", cfg)
	FrontendCfg = &cfg
	return nil
}
