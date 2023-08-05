package logger

import (
	"io"
	"os"
	"sync/atomic"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


type LogLevel int32

const (
	OffLevel LogLevel = iota
	InfoLevel
	DebugLevel
)

var logLevel = int32(InfoLevel)

func SetLogLevel(text string) error {
	level, err := parseLogLevel(text)
	if err != nil {
		return err
	}
	atomic.StoreInt32(&logLevel, int32(level))
	return nil
}

func parseLogLevel(text string) (LogLevel, error) {
	switch text {
	case "off":
		return OffLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	default:
		return OffLevel, errors.Errorf("unknown log level %s", text)			
	}
}

func NewLoggerWithName(name string) logr.Logger {
	return newLogger(LogLevel(logLevel)).WithName(name)
}

func newLogger(level LogLevel) logr.Logger {
	return newLoggerTo(os.Stdout, level)
}

func newLoggerTo(destWritter io.Writer, level LogLevel) logr.Logger {
	return zapr.NewLogger(newZapLoggerTo(destWritter, level))
}

func newZapLoggerTo(destWriter io.Writer, level LogLevel, opts ...zap.Option) *zap.Logger {
	var lvl zap.AtomicLevel
	switch level {
	case OffLevel:
		return zap.NewNop()
	case DebugLevel:
		lvl = zap.NewAtomicLevelAt(-10)
		opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	default:
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)		
	}

	encCfg := zap.NewDevelopmentEncoderConfig()
	enc := zapcore.NewConsoleEncoder(encCfg)
	sink := zapcore.AddSync(destWriter)
	opts = append(opts, zap.AddCallerSkip(1), zap.ErrorOutput(sink))
	return zap.New(zapcore.NewCore(enc, sink, lvl)).WithOptions(opts...)
}