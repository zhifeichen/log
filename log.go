package log

import (
	"context"
	"io"
	"os"
	"sync"
)

// Logger is the logging interface
type Logger interface {
	Trace(v ...any)
	Tracef(format string, v ...any)
	Debug(v ...any)
	Debugf(format string, v ...any)
	Info(v ...any)
	Infof(format string, v ...any)
	Warn(v ...any)
	Warnf(format string, v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Writer() io.Writer
	Flush() error
	Discard()
	ResumeWriter()
	SetContext(ctx context.Context) context.Context
	GetNewContext(ctx context.Context) context.Context
	WithContext(ctx context.Context) Logger
	WithTraceID(traceID uint32) Logger
}

// Level is a log level
type level int

// log level
const (
	LevelFatal level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

var levelString = map[level]string{
	LevelFatal: "Fatal",
	LevelError: "Error",
	LevelWarn:  "Warn",
	LevelInfo:  "Info",
	LevelDebug: "Debug",
	LevelTrace: "Trace",
}

var (
	defaultLogger Logger
	once          sync.Once
)

// getDefaultLogger returns the singleton logger, initializing it lazily if needed
func getDefaultLogger() Logger {
	once.Do(func() {
		defaultLogger = New(NewOptions(
			Filename("log.log"),
			Level("debug"),
		))
	})
	return defaultLogger
}

// Init init log
func Init(o Options) {
	defaultLogger = New(o)
}

// Flush flush any buffered log entries
func Flush() error {
	if defaultLogger != nil {
		return defaultLogger.Flush()
	}
	return nil
}

// Writer get log writer
func Writer() io.Writer {
	if defaultLogger != nil {
		return defaultLogger.Writer()
	}
	return os.Stderr
}

// Trace provides trace level logging
func Trace(v ...any) {
	getDefaultLogger().Trace(v...)
}

// Tracef provides trace level logging
func Tracef(format string, v ...any) {
	getDefaultLogger().Tracef(format, v...)
}

// Debug provides debug level logging
func Debug(v ...any) {
	getDefaultLogger().Debug(v...)
}

// Debugf provides debug level logging
func Debugf(format string, v ...any) {
	getDefaultLogger().Debugf(format, v...)
}

// Info provides info level logging
func Info(v ...any) {
	getDefaultLogger().Info(v...)
}

// Infof provides info level logging
func Infof(format string, v ...any) {
	getDefaultLogger().Infof(format, v...)
}

// Warn provides warn level logging
func Warn(v ...any) {
	getDefaultLogger().Warn(v...)
}

// Warnf provides warn level logging
func Warnf(format string, v ...any) {
	getDefaultLogger().Warnf(format, v...)
}

// Error provides error level logging
func Error(v ...any) {
	getDefaultLogger().Error(v...)
}

// Errorf provides error level logging
func Errorf(format string, v ...any) {
	getDefaultLogger().Errorf(format, v...)
}

// Fatal provides fatal level logging
func Fatal(v ...any) {
	getDefaultLogger().Fatal(v...)
	os.Exit(1)
}

// Fatalf provides fatal level logging
func Fatalf(format string, v ...any) {
	getDefaultLogger().Fatalf(format, v...)
	os.Exit(1)
}

// SetContext injects a random traceID into the context using the default logger
func SetContext(ctx context.Context) context.Context {
	return getDefaultLogger().SetContext(ctx)
}

// GetNewContext extracts or creates traceID using the default logger
func GetNewContext(ctx context.Context) context.Context {
	return getDefaultLogger().GetNewContext(ctx)
}

// WithContext returns a new Logger with traceID from context using the default logger
func WithContext(ctx context.Context) Logger {
	return getDefaultLogger().WithContext(ctx)
}

// WithTraceID returns a new Logger with the given traceID using the default logger
func WithTraceID(traceID uint32) Logger {
	return getDefaultLogger().WithTraceID(traceID)
}
