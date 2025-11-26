package logger

import (
	"log"
	"os"
)

var Log *Logger

type Logger struct {
	level string
}

func NewLogger(level string) *Logger {
	if level == "" {
		level = "info"
	}
	return &Logger{level: level}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level == "debug" {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level == "debug" || l.level == "info" {
		log.Printf("[INFO] "+format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

func init() {
	logLevel := os.Getenv("LOG_LEVEL")
	Log = NewLogger(logLevel)
}


