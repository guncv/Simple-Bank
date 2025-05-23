package worker

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Print(level zerolog.Level, args ...interface{}) {
	logger.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (l *Logger) Info(args ...interface{}) {
	l.Print(zerolog.InfoLevel, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Print(zerolog.ErrorLevel, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.Print(zerolog.DebugLevel, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Print(zerolog.FatalLevel, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.Print(zerolog.PanicLevel, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Print(zerolog.WarnLevel, args...)
}
