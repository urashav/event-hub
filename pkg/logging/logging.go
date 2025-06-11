package logging

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
)

type Logger struct {
	*log.Logger
	level Level
}

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

func NewLogger(level Level) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		level:  level,
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.output(DEBUG, format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.output(INFO, format, v...)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WARN {
		l.output(WARN, format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.output(ERROR, format, v...)
	}
}

func (l *Logger) output(level Level, format string, v ...interface{}) {
	// Получаем информацию о файле и строке
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	// Используем только имя файла
	file = path.Base(file)

	// Форматируем сообщение
	msg := fmt.Sprintf(format, v...)

	// Выводим лог в формате: [LEVEL] file:line message
	l.Printf("[%s] %s:%d %s", levelNames[level], file, line, msg)
}

// Глобальный логгер по умолчанию
var defaultLogger = NewLogger(INFO)

// Глобальные функции для удобства использования
func Debug(format string, v ...interface{}) {
	defaultLogger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	defaultLogger.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	defaultLogger.Error(format, v...)
}

// Функция для изменения уровня логирования глобального логгера
func SetLevel(level Level) {
	defaultLogger.level = level
}
