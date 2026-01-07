package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func NewLogger() Logger {
	writer := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}

	logger := zerolog.New(writer).
		With().
		Timestamp().
		Caller().
		Logger()

	return Logger{logger}
}

func (l *Logger) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *Logger) Error(err error, msg string) {
	l.Logger.Err(err).Msg(msg)
}

func (l *Logger) Debug(msg string) {
	l.Logger.Debug().Msg(msg)
}
