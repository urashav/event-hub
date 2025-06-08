package logging

import (
	"log"
)

type LoggerInterface interface {
	Debug(message string)
	Info(message string)
	Error(message string)
}

const (
	DEBUG = iota
	INFO
	ERROR
)

type Logger struct {
	Level int
	Layer string
}

func getLevel(level int) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case ERROR:
		return "ERROR"
	default:
		log.Printf("ERROR: Unknown log level: %d", level)
		return "UNKNOWN"
	}
}

func NewLogger(level int, layer string) *Logger {
	return &Logger{
		Level: level,
		Layer: layer,
	}
}
func (l *Logger) Debug(message string) {
	if l.Level == DEBUG {
		log.Printf("DEBUG: ["+l.Layer+"] %v", message)
	}
}

func (l *Logger) Info(message string) {
	if l.Level <= INFO {
		log.Printf("INFO: ["+l.Layer+"] %s ", message)
	}
}

func (l *Logger) Error(message string) {
	if l.Level <= ERROR {
		log.Printf("ERROR: [" + l.Layer + "] " + message)
	}
}
