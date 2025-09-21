package logger

import (
	"log"
	"os"
	"sync"
)

var (
	once   sync.Once
	logger *log.Logger
)

// GetLogger রিটার্ন করবে singleton logger
func GetLogger() *log.Logger {
	once.Do(func() {
		// ensure logs folder exists
		if _, err := os.Stat("storage/logs"); os.IsNotExist(err) {
			os.Mkdir("storage/logs", os.ModePerm)
		}

		file, err := os.OpenFile("storage/logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}

		logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	})
	return logger
}

// Laravel style wrapper functions
func Info(v ...any) {
	GetLogger().SetPrefix("INFO: ")
	GetLogger().Println(v...)
}

func Error(v ...any) {
	GetLogger().SetPrefix("ERROR: ")
	GetLogger().Println(v...)
}

func Warn(v ...any) {
	GetLogger().SetPrefix("WARN: ")
	GetLogger().Println(v...)
}
