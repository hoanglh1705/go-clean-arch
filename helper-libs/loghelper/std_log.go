package loghelper

import (
	"context"
	"log"

	"go.uber.org/zap"
)

var (
	StdLogger *stdLogger = &stdLogger{}
)

type debugWriter struct {
	appEnv string
}

func (s *debugWriter) Write(p []byte) (n int, err error) {
	if s.appEnv == "production" {
		zap.S().Debugf(string(p))
	} else if s.appEnv == "local" {
		log.Print(string(p))
	} else {
		log.Print(string(p))
	}

	return len(p), nil
}

type infoWriter struct {
	appEnv string
}

func (s *infoWriter) Write(p []byte) (n int, err error) {
	if s.appEnv == "production" {
		zap.S().Infof(string(p))
	} else if s.appEnv == "local" {
		log.Print(string(p))
	} else {
		log.Print(string(p))
	}

	return len(p), nil
}

type warningWriter struct {
	appEnv string
}

func (s *warningWriter) Write(p []byte) (n int, err error) {
	if s.appEnv == "production" {
		zap.S().Warnf(string(p))
	} else if s.appEnv == "local" {
		log.Print(string(p))
	} else {
		log.Print(string(p))
	}

	return len(p), nil
}

type errorWriter struct {
	appEnv string
}

func (s *errorWriter) Write(p []byte) (n int, err error) {
	if s.appEnv == "production" {
		zap.S().Errorf(string(p))
	} else if s.appEnv == "local" {
		log.Print(string(p))
	} else {
		log.Print(string(p))
	}

	return len(p), nil
}

func InitStdLogger(appEnv string) error {
	StdLogger = newStdLogger(appEnv)

	return nil
}

// func (s *Service) Write(p []byte) (n int, err error) {
// 	if s.appEnv == "production" {
// 		s.lokiService.Infof(string(p))
// 	} else if s.appEnv == "local" {
// 		log.Printf(string(p))
// 	}

// 	return len(p), nil
// }

type stdLogger struct {
	appEnv        string
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

func newStdLogger(appEnv string) *stdLogger {
	s := &stdLogger{
		appEnv:        appEnv,
		debugLogger:   log.New(&debugWriter{appEnv: appEnv}, "DEBUG: ", log.Ldate|log.Ltime|log.Llongfile),
		infoLogger:    log.New(&infoWriter{appEnv: appEnv}, "INFO: ", log.Ldate|log.Ltime|log.Llongfile),
		warningLogger: log.New(&warningWriter{appEnv: appEnv}, "WARNING: ", log.Ldate|log.Ltime|log.Llongfile),
		errorLogger:   log.New(&errorWriter{appEnv: appEnv}, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile),
	}

	return s
}

// Debugf uses fmt.Sprintf to log a templated message.
func (s *stdLogger) Debugf(template string, args ...interface{}) {
	s.debugLogger.Printf(template, args, nil)
}

// Infof uses fmt.Sprintf to log a templated message.
func (s *stdLogger) Infof(template string, args ...interface{}) {
	s.infoLogger.Printf(template, args, nil)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (s *stdLogger) Warnf(template string, args ...interface{}) {
	s.warningLogger.Printf(template, args, nil)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (s *stdLogger) Errorf(template string, args ...interface{}) {
	s.errorLogger.Printf(template, args, nil)
}

func (s *stdLogger) DebugfContext(ctx context.Context, template string, args ...interface{}) {
	s.debugLogger.Printf(template, args, nil)
}

func (s *stdLogger) InfofContext(ctx context.Context, template string, args ...interface{}) {
	s.infoLogger.Printf(template, args, nil)
}

func (s *stdLogger) WarnfContext(ctx context.Context, template string, args ...interface{}) {
	s.warningLogger.Printf(template, args, nil)
}

func (s *stdLogger) ErrorfContext(ctx context.Context, template string, args ...interface{}) {
	s.errorLogger.Printf(template, args, nil)
}
