package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"
	"strings"
)

var logger *logging

type logging struct {
	*zap.Logger
	verbose  int
	appid    string
	disabled bool
}

func loadSyncer(cfg *Config) zapcore.WriteSyncer {
	if cfg.EndPoint != "" {
		if cfg.Protocol == "" {
			cfg.Protocol = UDP
		}
		return NewAsyncNetWriter(cfg.EndPoint, cfg.Protocol).WriterSyncer()
	} else if cfg.LogFile != "" {
		file, err := os.Create(cfg.LogFile)
		if err != nil {
			panic(err)
		}
		return file
	}
	return os.Stdout
}

func defaultJsonEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		TimeKey:     "ts",
		EncodeTime:  zapcore.EpochTimeEncoder,
	})
}

func Initialize(cfg *Config) {
	if cfg == nil {
		Stdout()
		return
	}
	le := zapcore.InfoLevel
	if cfg.Debug {
		le = zapcore.DebugLevel
	}
	logger = &logging{
		Logger:   zap.New(zapcore.NewCore(defaultJsonEncoder(), loadSyncer(cfg), le)),
		verbose:  cfg.Verbose,
		appid:    cfg.AppId,
		disabled: cfg.Disable,
	}
}

func Stdout() {
	logger = &logging{
		Logger:  zap.New(zapcore.NewCore(defaultJsonEncoder(), os.Stdout, zapcore.DebugLevel)),
		verbose: 1,
	}
}

type Fields []zap.Field

func (f Fields) AppendHeader(depth int) Fields {
	_, file, line, ok := runtime.Caller(3 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	f = append(f, zap.String("file", file), zap.Int("line", line))
	return f
}

func (l *logging) log(verbose int, logF func(string, ...zap.Field), msg string, fields ...zap.Field) {
	if l.appid != "" {
		fields = append(fields, zap.String("appid", l.appid))
	}
	if l.verbose >= verbose {
		logF(msg, fields...)
	}
}

func verbose(level, depth int) Wrapper {
	if logger.disabled {
		return &voidWrapper{}
	}
	lw := &loggerWrapper{
		fields:  []zap.Field{},
		verbose: level,
	}
	lw.With(Fields([]zap.Field{zap.Int("verbose", level)}).AppendHeader(depth)...)
	return lw
}

func V(level int) Wrapper {
	return verbose(level, 0)
}

func With(fields ...zap.Field) Wrapper {
	return verbose(0, 0).With(fields...)
}

func Debugf(msg string, args ...interface{}) {
	verbose(0, 0).Debugf(msg, args...)
}

func Infof(msg string, args ...interface{}) {
	verbose(0, 0).Infof(msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	verbose(0, 0).Warnf(msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	verbose(0, 0).Errorf(msg, args...)
}

func DPanicf(msg string, args ...interface{}) {
	verbose(0, 0).DPanicf(msg, args...)
}

func Panicf(msg string, args ...interface{}) {
	verbose(0, 0).Panicf(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	verbose(0, 0).Fatalf(msg, args...)
}

func Debug(msg string, fields ...zap.Field) {
	verbose(0, 0).Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	verbose(0, 0).Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	verbose(0, 0).Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	verbose(0, 0).Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	verbose(0, 0).DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	verbose(0, 0).Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	verbose(0, 0).Fatal(msg, fields...)
}
