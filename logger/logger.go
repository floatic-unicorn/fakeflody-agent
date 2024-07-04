package logger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

var logger zerolog.Logger
var wlogger zerolog.Logger

func NewLogger(level string) {
	var logLevel zerolog.Level

	switch level {
	case "trace":
		logLevel = zerolog.TraceLevel
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(logLevel).With().Timestamp(). /*CallerWithSkipFrameCount(3).*/ Logger()
}

func NewFileLogger(level string, filePath string) error {
	var logLevel zerolog.Level

	switch level {
	case "trace":
		logLevel = zerolog.TraceLevel
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	wlogger = zerolog.New(
		file,
	).Level(logLevel).With().Timestamp().Logger()

	return nil
}

func Info(msg string) {
	logger.Info().Msg(msg)
}

func Error(msg string) {
	logger.Error().Msg(msg)
}

func Fatal(msg string) {
	logger.Fatal().Msg(msg)
}

func Debug(msg string) {
	logger.Debug().Msg(msg)
}

func Warn(msg string) {
	logger.Warn().Msg(msg)
}

func Trace(msg string) {
	logger.Trace().Msg(msg)
}

func Panic(msg string) {
	logger.Panic().Msg(msg)
}

func Infof(format string, v ...interface{}) {
	logger.Info().Msgf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	logger.Error().Msgf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatal().Msgf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debug().Msgf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	logger.Warn().Msgf(format, v...)
}

func Tracef(format string, v ...interface{}) {
	logger.Trace().Msgf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	logger.Panic().Msgf(format, v...)
}

func WInfo(msg string) {
	wlogger.Info().Msg(msg)
}

func WError(msg string) {
	wlogger.Error().Msg(msg)
}

func WFatal(msg string) {
	wlogger.Fatal().Msg(msg)
}

func WDebug(msg string) {
	wlogger.Debug().Msg(msg)
}

func WWarn(msg string) {
	wlogger.Warn().Msg(msg)
}

func WInfof(format string, v ...interface{}) {
	wlogger.Info().Msgf(format, v...)
}

func WWarnf(format string, v ...interface{}) {
	wlogger.Warn().Msgf(format, v...)
}

func WErrorf(msg string, v ...interface{}) {
	wlogger.Error().Msgf(msg, v)
}

func WFatalf(msg string, v ...interface{}) {
	wlogger.Fatal().Msgf(msg, v)
}

func WDebugf(msg string, v ...interface{}) {
	wlogger.Debug().Msgf(msg, v)
}
