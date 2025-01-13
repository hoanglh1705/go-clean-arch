package loghelper

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

var (
	SlogLogger *slogLogger = &slogLogger{}
)

func InitSlogTextLogger(appEnv string) error {
	SlogLogger = newSlogLogger(appEnv, "text")
	return nil
}

func InitSlogJSONLogger(appEnv string) error {
	SlogLogger = newSlogLogger(appEnv, "json")
	return nil
}

type slogLogger struct {
	slogLog *slog.Logger
}

func newSlogLogger(appEnv string, logFormat string) *slogLogger {
	slogLog := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if logFormat == "json" {
		slogLog = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return &slogLogger{
		slogLog: slogLog,
	}
}

// Debugf uses fmt.Sprintf to log a templated message.
func (s *slogLogger) Debugf(template string, args ...interface{}) {
	s.slogLog.Debug(fmt.Sprintf(template, args...))
}

// Infof uses fmt.Sprintf to log a templated message.
func (s *slogLogger) Infof(template string, args ...interface{}) {
	s.slogLog.Info(fmt.Sprintf(template, args...))
}

// Warnf uses fmt.Sprintf to log a templated message.
func (s *slogLogger) Warnf(template string, args ...interface{}) {
	s.slogLog.Warn(fmt.Sprintf(template, args...))
}

// Errorf uses fmt.Sprintf to log a templated message.
func (s *slogLogger) Errorf(template string, args ...interface{}) {
	s.slogLog.Error(fmt.Sprintf(template, args...))
}

func (s *slogLogger) DebugfContext(ctx context.Context, template string, args ...interface{}) {
	s.slogLog.Debug(fmt.Sprintf(template, args...), "traceId", getTraceIdFromContext(ctx))
}

func (s *slogLogger) InfofContext(ctx context.Context, template string, args ...interface{}) {
	s.slogLog.Info(fmt.Sprintf(template, args...), "traceId", getTraceIdFromContext(ctx))
}

func (s *slogLogger) WarnfContext(ctx context.Context, template string, args ...interface{}) {
	s.slogLog.Warn(fmt.Sprintf(template, args...), "traceId", getTraceIdFromContext(ctx))
}

func (s *slogLogger) ErrorfContext(ctx context.Context, template string, args ...interface{}) {
	s.slogLog.Error(fmt.Sprintf(template, args...), "traceId", getTraceIdFromContext(ctx))
}
