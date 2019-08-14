package log

import (
	"go.uber.org/zap"
)

type voidWrapper struct {
}

func (lw *voidWrapper) With(fields ...zap.Field) Wrapper {
	return lw
}

func (lw *voidWrapper) Debugf(msg string, args ...interface{}) {
	return
}

func (lw *voidWrapper) Infof(msg string, args ...interface{}) {
	return
}

func (lw *voidWrapper) Warnf(msg string, args ...interface{}) {
	return
}

func (lw *voidWrapper) Errorf(msg string, args ...interface{}) {
	return
}

func (lw *voidWrapper) DPanicf(msg string, args ...interface{}) {
	return
}

func (lw *voidWrapper) Panicf(msg string, args ...interface{}) {
	return
}

func (lw *voidWrapper) Fatalf(msg string, args ...interface{}) {
	return
}

func (lw *voidWrapper) Debug(msg string, fields ...zap.Field) {
	return
}

func (lw *voidWrapper) Info(msg string, fields ...zap.Field) {
	return
}

func (lw *voidWrapper) Warn(msg string, fields ...zap.Field) {
	return
}

func (lw *voidWrapper) Error(msg string, fields ...zap.Field) {
	return
}

func (lw *voidWrapper) DPanic(msg string, fields ...zap.Field) {
	return
}

func (lw *voidWrapper) Panic(msg string, fields ...zap.Field) {
	return
}

func (lw *voidWrapper) Fatal(msg string, fields ...zap.Field) {
	return
}
