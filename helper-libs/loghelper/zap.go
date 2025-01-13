package loghelper

import (
	"context"
	"go-clean-arch/helper-libs/commonhelper"
	"go-clean-arch/helper-libs/envhelper"
	"io"
	"os"

	"github.com/natefinch/lumberjack"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger   *zapLogger = &zapLogger{}
	DBLogger *zapLogger = &zapLogger{}
)

type zapLogger struct {
	*zap.SugaredLogger
	L *zap.Logger
}

// func InitZap(app, env string, maskFields map[string]string) error {
// 	logLevel := configLogLevel(env)
// 	encoderConfig := zapcore.EncoderConfig{
// 		MessageKey:   "message",
// 		LevelKey:     "level",
// 		EncodeLevel:  zapcore.CapitalLevelEncoder,
// 		TimeKey:      "time",
// 		EncodeTime:   zapcore.ISO8601TimeEncoder,
// 		CallerKey:    "caller",
// 		EncodeCaller: zapcore.ShortCallerEncoder,
// 		NameKey:      "app",
// 		EncodeName:   zapcore.FullNameEncoder,
// 	}
// 	// zap.RegisterEncoder("custom-json", func(ec zapcore.EncoderConfig) (zapcore.Encoder, error) {
// 	// 	return NewJSONEncoder(ec), nil
// 	// })
// 	cfg := zap.Config{
// 		// Encoding:         "custom-json",
// 		Encoding:         "json",
// 		Level:            zap.NewAtomicLevelAt(logLevel),
// 		OutputPaths:      []string{"stdout"},
// 		ErrorOutputPaths: []string{"stderr"},
// 		EncoderConfig:    encoderConfig,
// 	}

// 	l, err := cfg.Build()
// 	if err != nil {
// 		return err
// 	}
// 	l = l.Named(app)
// 	zap.ReplaceGlobals(l)
// 	Logger = &logger{l.Sugar(), l}

// 	return nil
// }

func InitZap(app, env string, maskingFields map[string]string) error {
	logLevel := configLogLevel(env)
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
		NameKey:      "app",
		EncodeName:   zapcore.FullNameEncoder,
	}

	syncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
	)

	newCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		syncer,
		zap.NewAtomicLevelAt(logLevel),
	)

	newLogger := zap.New(
		newCore,
		zap.AddCaller(),
	)
	defer func() {
		_ = newLogger.Sync()
	}()

	newLogger = newLogger.Named(app)
	zap.ReplaceGlobals(newLogger)
	Logger = &zapLogger{newLogger.Sugar(), newLogger}

	return nil
}

func InitZapWithWriters(app, env string, writers []io.Writer, maskingFields map[string]string) error {
	logLevel := configLogLevel(env)
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
		NameKey:      "app",
		EncodeName:   zapcore.FullNameEncoder,
	}

	writerSyncers := make([]zapcore.WriteSyncer, 0)
	for _, writer := range writers {
		writerSyncers = append(writerSyncers, zapcore.AddSync(writer))
	}
	syncer := zapcore.NewMultiWriteSyncer(
		writerSyncers...,
	)

	newCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		syncer,
		zap.NewAtomicLevelAt(logLevel),
	)

	newLogger := zap.New(
		newCore,
		zap.AddCaller(),
	)
	defer func() {
		_ = newLogger.Sync()
	}()

	newLogger = newLogger.Named(app)
	zap.ReplaceGlobals(newLogger)
	Logger = &zapLogger{newLogger.Sugar(), newLogger}

	return nil
}

func InitZapWithRotatingFile(app, env string, sqlOrmWriter io.Writer, maskFields map[string]string) error {
	logLevel := configLogLevel(env)
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
		NameKey:      "app",
		EncodeName:   zapcore.FullNameEncoder,
	}

	syncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(&lumberjack.Logger{
			Filename: "logs/log",
		}),
	)

	newCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		syncer,
		zap.NewAtomicLevelAt(logLevel),
	)

	newLogger := zap.New(
		newCore,
		zap.AddCaller(),
	)
	defer func() {
		_ = newLogger.Sync()
	}()

	newLogger = newLogger.Named(app)
	zap.ReplaceGlobals(newLogger)
	Logger = &zapLogger{newLogger.Sugar(), newLogger}

	return nil
}

func InitZapWithSql(app, env string, sqlOrmWriter io.Writer, maskFields map[string]string) error {
	logLevel := configLogLevel(env)
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
		NameKey:      "app",
		EncodeName:   zapcore.FullNameEncoder,
	}

	syncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(sqlOrmWriter),
	)

	newCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		syncer,
		zap.NewAtomicLevelAt(logLevel),
	)

	newLogger := zap.New(
		newCore,
		zap.AddCaller(),
	)
	defer func() {
		_ = newLogger.Sync()
	}()

	newLogger = newLogger.Named(app)
	DBLogger = &zapLogger{newLogger.Sugar(), newLogger}

	return nil
}

func (l *zapLogger) WithContext(ctx context.Context) *zap.SugaredLogger {
	return l.With(zap.Any("traceId", getTraceIdFromContext(ctx)))
}

func (l *zapLogger) WithCtx(ctx context.Context) *zap.SugaredLogger {
	return l.With(zap.Any("traceId", getTraceIdFromContext(ctx)))
}

func configLogLevel(defaultEnv string) zapcore.Level {
	env := os.Getenv(envhelper.ENVIRONMENT)
	if env == "" {
		env = defaultEnv
	}
	if env == "" {
		env = string(commonhelper.ENV__PRD)
	}

	var level zapcore.Level
	// level := zapcore.ErrorLevel
	// switch env {
	// case string(commonhelper.ENV__DEV):
	// 	level = zapcore.DebugLevel
	// case string(commonhelper.ENV__PRD):
	// 	level = zapcore.WarnLevel
	// default:
	// 	level = zapcore.WarnLevel
	// }

	logLevelEnv := os.Getenv(envhelper.LOG_LEVEL)
	switch logLevelEnv {
	case string(commonhelper.LOG_LEVEL__WARN):
		level = zapcore.WarnLevel
	case string(commonhelper.LOG_LEVEL__DEBUG):
		if env == string(commonhelper.ENV__PRD) {
			level = zapcore.InfoLevel
		} else {
			level = zapcore.DebugLevel
		}
	case string(commonhelper.LOG_LEVEL__INFO):
		level = zapcore.InfoLevel
	default:
		level = zapcore.InfoLevel
	}

	return level
}

func CreateFileRotatingWriter() io.Writer {
	return &lumberjack.Logger{
		Filename: "logs/log",
	}
}

func CreateStdoutWriter() io.Writer {
	return os.Stdout
}

func CreateStderrWriter() io.Writer {
	return os.Stderr
}
