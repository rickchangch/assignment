package log

import (
	"assignment-pe/internal/errs"
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	TimestampMilli string = "2006-01-02T15:04:05.000Z07:00"
)

type Level = string

const (
	LevelError Level = "error"
	LevelWarn  Level = "warn"
	LevelInfo  Level = "info"
	LevelDebug Level = "debug"
)

type Fields = logrus.Fields

type AppLogger interface {
	WithError(err error) AppLogger
	WithAppError(err errs.AppError) AppLogger
	WithField(key string, value interface{}) AppLogger
	WithFields(fields Fields) AppLogger
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
}

type appLogger struct {
	entry *logrus.Entry
}

type AppLoggerConfig struct {
	Level  Level
	Fields Fields
}

func NewLogger(config AppLoggerConfig) (AppLogger, error) {
	if config.Level == "" {
		config.Level = LevelDebug
	}

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, errors.New("invalid log level")
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: TimestampMilli})
	logger.SetOutput(os.Stderr)

	appLogger := &appLogger{}
	appLogger.entry = logrus.NewEntry(logger).WithFields(config.Fields)

	return appLogger, nil
}

func (al *appLogger) WithError(err error) AppLogger {
	_, file, line, _ := runtime.Caller(1)
	return &appLogger{
		entry: al.entry.
			WithError(err).
			WithField("stack_trace", fmt.Sprintf("%s:%d", file, line)),
	}
}

func (al *appLogger) WithAppError(err errs.AppError) AppLogger {
	return &appLogger{
		entry: al.entry.
			WithError(err).
			WithFields(Fields{
				"app_code":    err.AppCode,
				"stack_trace": err.GetStackTrace(),
			}),
	}
}

func (al *appLogger) WithField(key string, value interface{}) AppLogger {
	return &appLogger{entry: al.entry.WithField(key, value)}
}

func (al *appLogger) WithFields(fields Fields) AppLogger {
	return &appLogger{entry: al.entry.WithFields(fields)}
}

func (al *appLogger) Debug(args ...interface{})                 { al.entry.Debug(args...) }
func (al *appLogger) Debugf(format string, args ...interface{}) { al.entry.Debugf(format, args...) }
func (al *appLogger) Info(args ...interface{})                  { al.entry.Info(args...) }
func (al *appLogger) Infof(format string, args ...interface{})  { al.entry.Infof(format, args...) }
func (al *appLogger) Warn(args ...interface{})                  { al.entry.Warn(args...) }
func (al *appLogger) Warnf(format string, args ...interface{})  { al.entry.Warnf(format, args...) }
func (al *appLogger) Error(args ...interface{})                 { al.entry.Error(args...) }
func (al *appLogger) Errorf(format string, args ...interface{}) { al.entry.Errorf(format, args...) }
func (al *appLogger) Fatal(args ...interface{})                 { al.entry.Fatal(args...) }
func (al *appLogger) Fatalf(format string, args ...interface{}) { al.entry.Fatalf(format, args...) }
func (al *appLogger) Panic(args ...interface{})                 { al.entry.Panic(args...) }
func (al *appLogger) Panicf(format string, args ...interface{}) { al.entry.Panicf(format, args...) }
