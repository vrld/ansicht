package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogLevel string

const (
	LogLevelDebug   LogLevel = "DEBUG"
	LogLevelInfo    LogLevel = "INFO"
	logLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
)

type logger struct {
	mu   sync.Mutex
	file *os.File
}

var loggerInstance *logger

func Logger() *logger {
	if loggerInstance == nil {
		loggerInstance = &logger{}
	}
	return loggerInstance
}

func (l *logger) Initialize(logPath string) error {
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		l.file.Close()
	}

	l.file = file

	return nil
}

func (l *logger) Log(level LogLevel, message string) {
	if l.file == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	logLine := fmt.Sprintf("[%s|%s] %s\n", level, timestamp, message)

	if _, err := l.file.WriteString(logLine); err != nil {
		panic(fmt.Sprintf("failed to write to log file: %v", err))
	}

	l.file.Sync()
}

func (l *logger) Debug(message string) {
	l.Log(LogLevelDebug, message)
}

func (l *logger) Info(message string) {
	l.Log(LogLevelInfo, message)
}

func (l *logger) Warning(message string) {
	l.Log(logLevelWarning, message)
}

func (l *logger) Error(message string) {
	l.Log(LogLevelError, message)
}

func (l *logger) Close() {
	if l.file == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.file.Close()
	l.file = nil
}
