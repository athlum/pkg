package log

import (
	"go.uber.org/zap"
)

type Wrapper interface {
	With(fields ...zap.Field) Wrapper
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	DPanicf(msg string, args ...interface{})
	Panicf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}
