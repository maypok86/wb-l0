package logger

import (
	"go.uber.org/zap"
)

var instance *zap.SugaredLogger

func Debug(args ...interface{}) {
	instance.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	instance.Debugf(template, args...)
}

func Info(args ...interface{}) {
	instance.Info(args...)
}

func Infof(template string, args ...interface{}) {
	instance.Infof(template, args...)
}

func Warn(args ...interface{}) {
	instance.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	instance.Warnf(template, args...)
}

func Error(args ...interface{}) {
	instance.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	instance.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	instance.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	instance.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	instance.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	instance.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	instance.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	instance.Fatalf(template, args...)
}
