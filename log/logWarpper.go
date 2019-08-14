package log

import (
	"fmt"
	"go.uber.org/zap"
)

type loggerWrapper struct {
	fields  []zap.Field
	verbose int
}

func (lw *loggerWrapper) With(fields ...zap.Field) Wrapper {
	lw.fields = append(fields, lw.fields...)
	return lw
}

func (lw *loggerWrapper) logMessage(msg string, args ...interface{}) string {
	return fmt.Sprintf(msg, args...)
}

func (lw *loggerWrapper) Debugf(msg string, args ...interface{}) {
	logger.log(lw.verbose, logger.Debug, lw.logMessage(msg, args...), lw.fields...)
}

func (lw *loggerWrapper) Infof(msg string, args ...interface{}) {
	logger.log(lw.verbose, logger.Info, lw.logMessage(msg, args...), lw.fields...)
}

func (lw *loggerWrapper) Warnf(msg string, args ...interface{}) {
	logger.log(lw.verbose, logger.Warn, lw.logMessage(msg, args...), lw.fields...)
}

func (lw *loggerWrapper) Errorf(msg string, args ...interface{}) {
	logger.log(lw.verbose, logger.Error, lw.logMessage(msg, args...), lw.fields...)
}

func (lw *loggerWrapper) DPanicf(msg string, args ...interface{}) {
	logger.log(lw.verbose, logger.DPanic, lw.logMessage(msg, args...), lw.fields...)
}

func (lw *loggerWrapper) Panicf(msg string, args ...interface{}) {
	logger.log(lw.verbose, logger.Panic, lw.logMessage(msg, args...), lw.fields...)
}

func (lw *loggerWrapper) Fatalf(msg string, args ...interface{}) {
	logger.log(lw.verbose, logger.Fatal, lw.logMessage(msg, args...), lw.fields...)
}

func (lw *loggerWrapper) log(logF func(string, ...zap.Field), msg string, fields ...zap.Field) {
	lw.With(fields...)
	logger.log(lw.verbose, logF, msg, lw.fields...)
}

func (lw *loggerWrapper) Debug(msg string, fields ...zap.Field) {
	lw.log(logger.Debug, msg, fields...)
}

func (lw *loggerWrapper) Info(msg string, fields ...zap.Field) {
	lw.log(logger.Info, msg, fields...)
}

func (lw *loggerWrapper) Warn(msg string, fields ...zap.Field) {
	lw.log(logger.Warn, msg, fields...)
}

func (lw *loggerWrapper) Error(msg string, fields ...zap.Field) {
	lw.log(logger.Error, msg, fields...)
}

func (lw *loggerWrapper) DPanic(msg string, fields ...zap.Field) {
	lw.log(logger.DPanic, msg, fields...)
}

func (lw *loggerWrapper) Panic(msg string, fields ...zap.Field) {
	lw.log(logger.Panic, msg, fields...)
}

func (lw *loggerWrapper) Fatal(msg string, fields ...zap.Field) {
	lw.log(logger.Fatal, msg, fields...)
}
