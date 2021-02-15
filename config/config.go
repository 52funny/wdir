package config

import (
	"os"
	"strconv"

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
	// ShowHiddenFiles is show Hidden Files
	ShowHiddenFiles bool
)

// ReadConfig is read config function
func ReadConfig(configName string) error {
	envWdirDocker := os.Getenv("WDIR_DOCKER")
	isDocker, _ := strconv.ParseBool(envWdirDocker)
	if isDocker {
		Port = os.Getenv("PORT")
		Template = os.Getenv("TEMPLATE")
		Path = os.Getenv("FILEPATH")
		LogPath = os.Getenv("LOGPATH")
		ShowHiddenFiles, _ = strconv.ParseBool(os.Getenv("SHOWHIDDENFILES"))
		return nil
	}
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
	LogPath = viper.GetString("config.log_path")
	ShowHiddenFiles = viper.GetBool("config.show_hidden_files")
	return nil
}
