package config

import (
	"strconv"
	"syscall"

	"github.com/jinzhu/configor"
)

var Config = struct {
	Port            string `default:"80" yaml:"port"`
	Path            string `default:"." yaml:"path"`
	LogPath         string `default:"log" yaml:"log_path"`
	ShowHiddenFiles bool   `default:"false" yaml:"show_hidden_files"`
}{}

// ReadConfig is read config function
func ReadConfig(configName, flagPort, flagPath string) error {
	// first read from file
	configor.Load(&Config, configName)

	Config.Port = flagPort
	Config.Path = flagPath

	port, ok := syscall.Getenv("PORT")
	if ok {
		Config.Port = port
	}
	path, ok := syscall.Getenv("FILEPATH")
	if ok {
		Config.Path = path
	}
	logPath, ok := syscall.Getenv("LOGPATH")
	if ok {
		Config.LogPath = logPath
	}

	show, ok := syscall.Getenv("SHOWHIDDENFILES")
	if ok {
		Config.ShowHiddenFiles, _ = strconv.ParseBool(show)
	}
	return nil
}
