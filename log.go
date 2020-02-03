package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Level is a log level
type Level int

// log level
const (
	LevelFatal Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

var levelString = map[Level]string{
	LevelFatal: "Fatal",
	LevelError: "Error",
	LevelWarn: "Warn",
	LevelInfo: "Info",
	LevelDebug: "Debug",
	LevelTrace: "Trace",
}

var (
	level = LevelTrace
)

// Init init log
func Init(l Level, path string) {
	level = l
	log.SetOutput(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    10, // MB
		MaxBackups: 20,
		MaxAge:     30,
	})
	log.SetFlags(log.Ldate | log.Ltime)
}

// Log makes use of log
func Log(file string, line int, l Level, v ...interface{}) {
	fl := fmt.Sprintf("[%5s] %s:%d: ", levelString[l], file, line)
	lv := fmt.Sprint(v...)
	log.Println(fl, lv)
}

// Logf makes use of log
func Logf(file string, line int, l Level, format string, v ...interface{}) {
	fl := fmt.Sprintf("[%5s] %s:%d: ", levelString[l], file, line)
	lv := fmt.Sprintf(format, v...)
	log.Print(fl, lv)
}

// WithLevel logs with the level specified
func WithLevel(file string, line int, l Level, v ...interface{}) {
	if l > level {
		return
	}
	Log(file, line, l, v...)
}

// WithLevelf logs with the level specified
func WithLevelf(file string, line int, l Level, format string, v ...interface{}) {
	if l > level {
		return
	}
	Logf(file, line, l, format, v...)
}

// Trace provides trace level logging
func Trace(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevel(filepath.Base(file), line, LevelTrace, v...)
}

// Tracef provides trace level logging
func Tracef(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevelf(filepath.Base(file), line, LevelTrace, format, v...)
}

// Debug provides debug level logging
func Debug(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevel(filepath.Base(file), line, LevelDebug, v...)
}

// Debugf provides debug level logging
func Debugf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevelf(filepath.Base(file), line, LevelDebug, format, v...)
}

// Info provides info level logging
func Info(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevel(filepath.Base(file), line, LevelInfo, v...)
}

// Infof provides info level logging
func Infof(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevelf(filepath.Base(file), line, LevelInfo, format, v...)
}

// Warn provides warn level logging
func Warn(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevel(filepath.Base(file), line, LevelWarn, v...)
}

// Warnf provides warn level logging
func Warnf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevelf(filepath.Base(file), line, LevelWarn, format, v...)
}

// Error provides error level logging
func Error(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevel(filepath.Base(file), line, LevelError, v...)
}

// Errorf provides error level logging
func Errorf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevelf(filepath.Base(file), line, LevelError, format, v...)
}

// Fatal provides fatal level logging
func Fatal(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevel(filepath.Base(file), line, LevelFatal, v...)
	os.Exit(1)
}

// Fatalf provides fatal level logging
func Fatalf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WithLevelf(filepath.Base(file), line, LevelFatal, format, v...)
	os.Exit(1)
}

// Writer get log writer
func Writer() io.Writer {
	return log.Writer()
}
