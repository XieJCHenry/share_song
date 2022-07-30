package logger

import (
	"share_song/global"

	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Sugared() *zap.SugaredLogger {
	return l.logger
}

func (l *Logger) Usable() bool {
	return l.logger != nil
}

func (l *Logger) Key() string {
	return global.KeyLogger
}
