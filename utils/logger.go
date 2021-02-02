package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	Log *log.Logger
)

// InitLogger is get a custom logger
func InitLogger(path string) {
	if !PathExists(path) {
		os.MkdirAll(path, 0773)
	}
	newPath := filepath.Join(path, time.Now().Format("2006-01-02"+".log"))
	f, err := os.OpenFile(newPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0773)
	if err != nil {
		panic(err)
	}
	writers := []io.Writer{
		f,
		os.Stdout,
	}
	Log = log.New(io.MultiWriter(writers...), "[wdir] ", log.Ldate|log.Ltime|log.Lshortfile)
}
