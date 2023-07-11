package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Level string
}

var logger *logrus.Logger

func init() {
	logLevel := os.Getenv("LOGGER_LEVEL")
	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Errorf("wrong log level, %v", logLevel)
		}
		logrus.SetLevel(level)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:  time.RFC3339,
		DisableTimestamp: false,
	})
	logger = logrus.New()
}

// func Initialize(config *Config) error {
// 	if config.Level != "" {
// 		level, err := logrus.ParseLevel(config.Level)
// 		if err != nil {
// 			return errors.WithStack(err)
// 		}
// 		logrus.SetLevel(level)
// 	}

// 	logrus.SetFormatter(&logrus.JSONFormatter{
// 		TimestampFormat:  time.RFC3339,
// 		DisableTimestamp: false,
// 	})

// 	return nil
// }

func WithName(name string) *logrus.Logger {
	return logger.WithField("package", name).Logger
}

// func NewLogger() *logrus.Logger {
// 	return logrus.New()
// }
