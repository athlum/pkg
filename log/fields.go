package log

import (
	"go.uber.org/zap"
)

func Type(logType string) zap.Field {
	return zap.String("logType", logType)
}

func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}
