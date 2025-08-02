package worker

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type customLogger struct {
}

func NewLogger() *customLogger {
	return &customLogger{}
}

func (l *customLogger) prefixPrint(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (l *customLogger) Debug(args ...interface{}) {
	l.prefixPrint(zerolog.DebugLevel, args...)

}

func (l *customLogger) Info(args ...interface{}) {
	l.prefixPrint(zerolog.InfoLevel, args...)
}

func (l *customLogger) Warn(args ...interface{}) {
	l.prefixPrint(zerolog.WarnLevel, args...)
}

func (l *customLogger) Error(args ...interface{}) {
	l.prefixPrint(zerolog.ErrorLevel, args...)
}

func (l *customLogger) Fatal(args ...interface{}) {
	l.prefixPrint(zerolog.FatalLevel, args...)
}
