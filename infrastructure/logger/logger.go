package logger

import (
	"log/slog"
	"os"
	"sync"
)

// Logger wraps slog.Logger to provide structured logging
type Logger struct {
	logger *slog.Logger
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{logger: l.logger.With(args...)}
}

var (
	logger *Logger
	once   sync.Once
)

// GetLogger returns the singleton logger instance
func GetLogger() *Logger {
	once.Do(SetupLogger)
	return logger
}

// SetupLogger initializes the logger
func SetupLogger() {
	slogInstance := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
	logger = &Logger{logger: slogInstance}
}

// SetLogger sets a custom logger implementation
func SetLogger(l *Logger) {
	logger = l
}
