package logger

import (
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Level string
	//SentryEnabled bool
}

const (
	sentryDSN = ""
)

func Initialize(config *Config) error {
	if config.Level != "" {
		level, err := logrus.ParseLevel(config.Level)
		if err != nil {
			return errors.WithStack(err)
		}
		logrus.SetLevel(level)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:  time.RFC3339,
		DisableTimestamp: false,
	})

	//if config.SentryEnabled {
	//	hook, err := logrus_sentry.NewAsyncSentryHook(sentryDSN, []logrus.Level{
	//		logrus.PanicLevel,
	//		logrus.FatalLevel,
	//		logrus.ErrorLevel,
	//	})
	//
	//	if err != nil {
	//		panic("failed to setup raven!")
	//	}
	//
	//	hook.StacktraceConfiguration.Enable = true
	//
	//	logrus.AddHook(hook)
	//}

	return nil
}

func NewLogger() *logrus.Logger {
	return logrus.New()
}
