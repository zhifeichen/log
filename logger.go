package log

import (
	"context"
	"io"
	"math/rand/v2"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const traceIDKey = "traceID"

// loggerImpl implements the Logger interface using zap
type loggerImpl struct {
	zapLogger *zap.SugaredLogger
	discarded bool
	mu        sync.RWMutex
}

// New creates a new Logger with the given options
func New(o Options) Logger {
	conf := o.toZapConfig()

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(conf.level)); err != nil {
		level = zapcore.DebugLevel
	}

	var w io.Writer = &lumberjack.Logger{
		Filename:   conf.filename,
		MaxSize:    conf.maxSize,
		MaxBackups: conf.maxBackups,
		MaxAge:     conf.maxAge,
		LocalTime:  true,
	}
	if len(conf.writers) > 0 {
		all := make([]io.Writer, 0, len(conf.writers)+1)
		all = append(all, conf.writers...)
		all = append(all, w)
		w = io.MultiWriter(all...)
	}

	fileSyncer := zapcore.AddSync(w)

	encConf := zap.NewProductionEncoderConfig()
	encConf.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000")
	encConf.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encConf)

	core := zapcore.NewCore(encoder, fileSyncer, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	return &loggerImpl{
		zapLogger: logger.Sugar(),
	}
}

// Discard discards all log output
func (l *loggerImpl) Discard() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.discarded = true
}

// ResumeWriter resumes log output after Discard
func (l *loggerImpl) ResumeWriter() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.discarded = false
}

// Flush flushes any buffered log entries
func (l *loggerImpl) Flush() error {
	return l.zapLogger.Sync()
}

func (l *loggerImpl) isDiscarded() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.discarded
}

// Trace provides trace level logging
func (l *loggerImpl) Trace(v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Debug(v...)
}

// Tracef provides trace level logging
func (l *loggerImpl) Tracef(format string, v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Debugf(format, v...)
}

// Debug provides debug level logging
func (l *loggerImpl) Debug(v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Debug(v...)
}

// Debugf provides debug level logging
func (l *loggerImpl) Debugf(format string, v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Debugf(format, v...)
}

// Info provides info level logging
func (l *loggerImpl) Info(v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Info(v...)
}

// Infof provides info level logging
func (l *loggerImpl) Infof(format string, v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Infof(format, v...)
}

// Warn provides warn level logging
func (l *loggerImpl) Warn(v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Warn(v...)
}

// Warnf provides warn level logging
func (l *loggerImpl) Warnf(format string, v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Warnf(format, v...)
}

// Error provides error level logging
func (l *loggerImpl) Error(v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Error(v...)
}

// Errorf provides error level logging
func (l *loggerImpl) Errorf(format string, v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Errorf(format, v...)
}

// Fatal provides fatal level logging
func (l *loggerImpl) Fatal(v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Fatal(v...)
}

// Fatalf provides fatal level logging
func (l *loggerImpl) Fatalf(format string, v ...any) {
	if l.isDiscarded() {
		return
	}
	l.zapLogger.Fatalf(format, v...)
}

// Write implements io.Writer
func (l *loggerImpl) Write(p []byte) (n int, err error) {
	if l.isDiscarded() {
		return len(p), nil
	}
	l.zapLogger.Info(string(p))
	return len(p), nil
}

// Writer returns an io.Writer that writes to this logger at Info level
func (l *loggerImpl) Writer() io.Writer {
	return l
}

// SetContext injects a random traceID into the context
func (l *loggerImpl) SetContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, traceIDKey, rand.Uint32())
}

// GetNewContext extracts traceID from existing context, or generates new one
func (l *loggerImpl) GetNewContext(ctx context.Context) context.Context {
	newCtx := context.Background()
	if ctx == nil {
		return context.WithValue(newCtx, traceIDKey, rand.Uint32())
	}
	traceID, ok := ctx.Value(traceIDKey).(uint32)
	if !ok {
		traceID = rand.Uint32()
	}
	return context.WithValue(newCtx, traceIDKey, traceID)
}

// WithContext returns a new Logger with traceID from context
func (l *loggerImpl) WithContext(ctx context.Context) Logger {
	if l.isDiscarded() {
		return l
	}
	if ctx == nil {
		return l
	}
	traceID, ok := ctx.Value(traceIDKey).(uint32)
	if !ok {
		return l
	}
	return &loggerImpl{
		zapLogger: l.zapLogger.With(zap.Uint32("traceID", traceID)),
	}
}

// WithTraceID returns a new Logger with the given traceID
func (l *loggerImpl) WithTraceID(traceID uint32) Logger {
	if l.isDiscarded() {
		return l
	}
	return &loggerImpl{
		zapLogger: l.zapLogger.With(zap.Uint32("traceID", traceID)),
	}
}
