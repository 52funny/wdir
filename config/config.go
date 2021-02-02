package config

import (
	"os"

	"github.com/spf13/viper"
)

var (
	// Port is server port
	Port string
	//Template is template path
	Template string
	// Path is dir path
	Path string
	// LogPath is log path
	LogPath string
)

// ReadConfig is read config function
func ReadConfig(configName string) error {
	viper.SetConfigName(configName)
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	viper.AddConfigPath(dir)
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	Port = viper.GetString("config.port")
	Template = viper.GetString("config.template")
	Path = viper.GetString("config.path")
	LogPath = viper.GetString("config.logpath")
	return nil
}
