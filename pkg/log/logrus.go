package log

import (
	"github.com/sirupsen/logrus"
)

// Logger abstruct interface for internal logging
type Logger interface {
	Error(msgs ...any)
	Warn(msgs ...any)
	Info(msgs ...any)
	Debug(msgs ...any)
	Trace(msgs ...any)
	Tracef(s string, msgs ...any)
	WithError(err error) Logger
	WithField(key string, value any) Logger
}

type logger struct {
	log *logrus.Entry
}

func (l *logger) Warn(msgs ...any) {
	l.log.Warn(msgs...)
}

func (l *logger) Tracef(s string, msgs ...any) {
	l.log.Tracef(s, msgs...)
}

func (l *logger) WithField(key string, value any) Logger {
	return NewLoggerFromEntry(l.log.WithField(key, value))
}

func (l *logger) WithError(err error) Logger {
	return NewLoggerFromEntry(l.log.WithError(err))
}

func (l *logger) Error(msgs ...any) {
	l.log.Error(msgs...)
}

func (l *logger) Info(msgs ...any) {
	l.log.Info(msgs...)
}

func (l *logger) Debug(msgs ...any) {
	l.log.Debug(msgs...)
}

func (l *logger) Trace(msgs ...any) {
	l.log.Trace(msgs...)
}

// NewLogger construct Logger from logrus.Logger
func NewLogger(log *logrus.Logger) Logger {
	return &logger{log: logrus.NewEntry(log)}
}

// NewLoggerFromEntry construct Logger from logrus.Entry
func NewLoggerFromEntry(log *logrus.Entry) Logger {
	return &logger{log: log}
}
