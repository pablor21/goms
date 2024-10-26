package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerConfig struct {
	Level int `json:"level" yaml:"level" MAPSTRUCTURE:"level" default:"1"`
}

func InitLogger(config LoggerConfig) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.Level(config.Level))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}).With().Timestamp().Logger()

}

func Logger() zerolog.Logger {
	return log.Logger
}

func Info() *zerolog.Event {
	return log.Info()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Trace() *zerolog.Event {
	return log.Trace()
}

func Panic() *zerolog.Event {
	return log.Panic()
}

func WithLevel(level zerolog.Level) *zerolog.Event {
	return log.WithLevel(level)
}
