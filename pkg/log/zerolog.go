package log

import (
	"github.com/rs/zerolog"
)

type zerologLogger struct {
	log *zerolog.Logger
}

func (z *zerologLogger) Error(msgs ...any) {
	log := z.log.Error()

	for _, msg := range msgs {
		log = log.Any("msg", msg)
	}

	log.Send()
}

func (z *zerologLogger) Warn(msgs ...any) {
	log := z.log.Warn()

	for _, msg := range msgs {
		log = log.Any("msg", msg)
	}

	log.Send()
}

func (z *zerologLogger) Info(msgs ...any) {
	log := z.log.Info()

	for _, msg := range msgs {
		log = log.Any("msg", msg)
	}

	log.Send()
}

func (z *zerologLogger) Debug(msgs ...any) {
	log := z.log.Debug()

	for _, msg := range msgs {
		log = log.Any("msg", msg)
	}

	log.Send()
}

func (z *zerologLogger) Trace(msgs ...any) {
	log := z.log.Trace()

	for _, msg := range msgs {
		log = log.Any("msg", msg)
	}

	log.Send()
}

func (z *zerologLogger) Tracef(s string, msgs ...any) {
	log := z.log.Trace()

	log.Msgf(s, msgs)
}

func (z *zerologLogger) WithError(err error) Logger {
	log := z.log.With().Err(err).Logger()
	return &zerologLogger{log: &log}
}

func (z *zerologLogger) WithField(key string, value any) Logger {
	log := z.log.With().Any(key, value).Logger()
	return &zerologLogger{log: &log}
}

func NewZerologLogger(logger *zerolog.Logger) Logger {
	return &zerologLogger{log: logger}
}
