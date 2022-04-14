package logger

import (
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type logger struct {
	log       zerolog.Logger
	writers   []io.Writer
	verbosity int
}

func Init(opts ...InitOption) zerolog.Logger {
	// prepare logger
	l := &logger{
		writers:   make([]io.Writer, 0),
		verbosity: 0,
	}

	// loop options
	for _, opt := range opts {
		opt(l)
	}

	// set logger
	l.log = log.Output(io.MultiWriter(l.writers...))

	// set global logger
	switch {
	case l.verbosity == 1:
		log.Logger = l.log.Level(zerolog.DebugLevel)
	case l.verbosity > 1:
		log.Logger = l.log.Level(zerolog.TraceLevel)
	default:
		log.Logger = l.log.Level(zerolog.InfoLevel)
	}

	return log.Logger
}

func Child(opts ...ChildOption) zerolog.Logger {
	// default
	l := log.With().
		Logger()

	// loop options
	for _, opt := range opts {
		opt(&l)
	}

	return l
}
