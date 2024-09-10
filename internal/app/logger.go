package app

import (
	"io"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func NewLogger(stdout io.Writer) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := zerolog.New(stdout).With().Timestamp().Logger()

	return logger
}
