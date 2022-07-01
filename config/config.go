package config

import (
	"os"
	"strconv"

	"github.com/jinzhu/configor"
)

var Config = struct {
	Port            string `default:"80" yaml:"port"`
	Path            string `default:"." yaml:"path"`
	LogPath         string `default:"log" yaml:"log_path"`
	ShowHiddenFiles bool   `default:"false" yaml:"show_hidden_files"`
}{}

// ReadConfig is read config function
func ReadConfig(configName string) error {
	envWdirDocker := os.Getenv("WDIR_DOCKER")
	isDocker, _ := strconv.ParseBool(envWdirDocker)
	// if is docker
	if isDocker {
		Config.Port = os.Getenv("PORT")
		Config.Path = os.Getenv("FILEPATH")
		Config.LogPath = os.Getenv("LOGPATH")
		Config.ShowHiddenFiles, _ = strconv.ParseBool(os.Getenv("SHOWHIDDENFILES"))
		return nil
	}
	configor.Load(&Config, configName)
	return nil
}
