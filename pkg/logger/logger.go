package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	})
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Error(msg string, err error) {
	log.Error().Err(err).Msg(msg)
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Infof(msg string, args ...interface{}) {
	log.Info().Msgf(msg, args...)
}

func Errorf(msg string, err error, args ...interface{}) {
	log.Error().Err(err).Msgf(msg, args...)
}
