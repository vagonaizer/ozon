package logging

import (
	"log"
	"os"
	"sync"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var levelPrefixes = map[Level]string{
	DEBUG: "[DEBUG] ",
	INFO:  "[INFO] ",
	WARN:  "[WARN] ",
	ERROR: "[ERROR] ",
}

// Logger структура
type Logger struct {
	mu    sync.Mutex
	level Level
	out   *log.Logger
}

var (
	instance *Logger
	once     sync.Once
)

func Init(level Level) {
	once.Do(func() {
		instance = &Logger{
			level: level,
			out:   log.New(os.Stdout, "", log.LstdFlags),
		}
	})
}

func GetLogger() *Logger {
	if instance == nil {
		Init(INFO)
	}
	return instance
}

func (l *Logger) logf(lvl Level, format string, args ...interface{}) {
	if lvl >= l.level {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.out.SetPrefix(levelPrefixes[lvl])
		l.out.Printf(format, args...)
	}
}

func (l *Logger) Info(format string, args ...interface{})  { l.logf(INFO, format, args...) }
func (l *Logger) Warn(format string, args ...interface{})  { l.logf(WARN, format, args...) }
func (l *Logger) Error(format string, args ...interface{}) { l.logf(ERROR, format, args...) }
func (l *Logger) Debug(format string, args ...interface{}) { l.logf(DEBUG, format, args...) }
