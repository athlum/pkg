package log

import (
	"go.uber.org/zap"
)

type StatslibLogger struct {
	loggerWrapper
}

func (sl *StatslibLogger) Printf(format string, v ...interface{}) {
	sl.Errorf(format, v...)
}

func NewStatslibLogger() *StatslibLogger {
	return &StatslibLogger{
		loggerWrapper: loggerWrapper{
			verbose: 1,
			fields: []zap.Field{
				Type("statslib"),
			},
		},
	}
}
