package logger

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Logger 日志实现;
var Logger = logrus.New()

func init() {
	Logger.AddHook(hook())
	Logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})
}

func hook() logrus.Hook {
	logrus.New()
	writer, err := rotatelogs.New("./logger"+".%Y%m%d%H",
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(time.Hour*24*365),
	)

	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
	}
	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05", PrettyPrint: false})
}

// Info info.
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// InfoF info format.
func InfoF(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// InfoM info fields map[string]interface{}.
func InfoM(fields logrus.Fields, args ...interface{}) {
	Logger.WithFields(fields).Info(args...)
}

// Warn warn.
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// WarnF warn format.
func WarnF(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// WarnM warn fields map[string]interface{}.
func WarnM(fields logrus.Fields, args ...interface{}) {
	Logger.WithFields(fields).Warn(args...)
}

// Error error.
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// ErrorF error format.
func ErrorF(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// ErrorM error fields map[string]interface{}.
func ErrorM(fields logrus.Fields, args ...interface{}) {
	Logger.WithFields(fields).Error(args...)
}

// Fatal fatal.
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// FatalF fatal format.
func FatalF(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// FatalM fatal fields map[string]interface{}.
func FatalM(fields logrus.Fields, args ...interface{}) {
	Logger.WithFields(fields).Fatal(args...)
}

// Panic panic.
func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

// PanicF panic format.
func PanicF(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

// PanicM panic fields map[string]interface{}.
func PanicM(fields logrus.Fields, args ...interface{}) {
	Logger.WithFields(fields).Panic(args...)
}
