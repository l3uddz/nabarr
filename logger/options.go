package logger

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

/* core logger options */

type InitOption func(*logger)

func WithConsole() InitOption {
	return func(l *logger) {
		l.writers = append(l.writers, zerolog.ConsoleWriter{
			TimeFormat: time.Stamp,
			Out:        os.Stderr,
		})
	}
}

func WithFile(filepath string) InitOption {
	return func(l *logger) {
		l.writers = append(l.writers, zerolog.ConsoleWriter{
			TimeFormat: time.Stamp,
			Out: &lumberjack.Logger{
				Filename:   filepath,
				MaxSize:    5,
				MaxAge:     14,
				MaxBackups: 5,
			},
			NoColor: true,
		})
	}
}

func WithVerbosity(verbosity int) InitOption {
	return func(l *logger) {
		l.verbosity = verbosity
	}
}

/* child logger options */

type ChildOption func(*zerolog.Logger)

func WithName(name string) ChildOption {
	return func(l *zerolog.Logger) {
		*l = l.With().
			Str("caller", name).
			Logger()
	}
}

func WithLevel(level string) ChildOption {
	return func(l *zerolog.Logger) {
		if level == "" {
			return
		}

		zl, err := zerolog.ParseLevel(level)
		if err != nil {
			return
		}

		*l = l.Level(zl)
	}
}
