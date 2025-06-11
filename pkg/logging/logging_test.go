package logging

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// Создаем буфер для перехвата вывода
	var buf bytes.Buffer

	// Создаем логгер с выводом в буфер
	logger := &Logger{
		Logger: log.New(&buf, "", log.Ldate|log.Ltime),
		level:  DEBUG,
	}

	tests := []struct {
		name    string
		level   Level
		logFunc func(string, ...interface{})
		message string
		want    string
	}{
		{
			name:    "debug message",
			level:   DEBUG,
			logFunc: logger.Debug,
			message: "test debug",
			want:    "[DEBUG]",
		},
		{
			name:    "info message",
			level:   INFO,
			logFunc: logger.Info,
			message: "test info",
			want:    "[INFO]",
		},
		{
			name:    "warn message",
			level:   WARN,
			logFunc: logger.Warn,
			message: "test warn",
			want:    "[WARN]",
		},
		{
			name:    "error message",
			level:   ERROR,
			logFunc: logger.Error,
			message: "test error",
			want:    "[ERROR]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.message)

			got := buf.String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("Logger output = %v, want %v", got, tt.want)
			}
			if !strings.Contains(got, tt.message) {
				t.Errorf("Logger output = %v, should contain message %v", got, tt.message)
			}
		})
	}
}

func TestLogLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		Logger: log.New(&buf, "", log.Ldate|log.Ltime),
		level:  INFO, // Устанавливаем уровень INFO
	}

	// Debug сообщение не должно появиться
	logger.Debug("debug message")
	if buf.String() != "" {
		t.Error("Debug message should not appear when level is INFO")
	}

	// Info сообщение должно появиться
	buf.Reset()
	logger.Info("info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Error("Info message should appear when level is INFO")
	}
}
