package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	LogFile *os.File
)

func InitLogger(logFilePath string) error {
	var err error
	LogFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	log.SetOutput(io.MultiWriter(LogFile, os.Stdout))
	return nil
}

func CleanupOldLogs(logDir string, maxAge int) {
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		log.Println("Error reading log directory:", err)
		return
	}

	for _, file := range files {
		if time.Since(file.ModTime()) > time.Duration(maxAge)*24*time.Hour {
			os.Remove(file.Name())
		}
	}
}
