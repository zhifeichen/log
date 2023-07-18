package log

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

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
	logLevel  = LevelTrace
	bufWriter *bufio.Writer
	debounced debounce
	logChn    chan string
)

// Init init log
func Init(o Options) {
	logLevel = o.level
	var w io.Writer = &lumberjack.Logger{
		Filename:   o.filename,
		MaxSize:    o.maxSize, // MB
		MaxBackups: o.maxBackups,
		MaxAge:     o.maxAge,
		LocalTime:  true,
	}
	if len(o.writers) > 0 {
		o.writers = append(o.writers, w)
		w = io.MultiWriter(o.writers...)
	}
	bufWriter = bufio.NewWriterSize(w, 64*1024)
	debounced = newDebouncer(time.Second, func() {
		bufWriter.Flush()
	})
	logChn = make(chan string, 50)
	go logLoop()
	log.SetOutput(bufWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

// Flush flush buf to file
func Flush() error {
	if bufWriter != nil {
		return bufWriter.Flush()
	}
	return nil
}

func logLoop() {
	for {
		debounced()
		l, ok := <-logChn
		if !ok {
			return
		}
		if bufWriter != nil && len(l) > bufWriter.Available() {
			bufWriter.Flush()
		}
		log.Print(l)
	}
}

// Log makes use of log
func Log(file string, line int, l level, v ...interface{}) {
	fl := fmt.Sprintf("[%5s] %s:%d:", levelString[l], file, line)
	lv := fmt.Sprint(v...)
	// log.Println(fl, lv)
	str := fmt.Sprintln(fl, lv)
	logChn <- str
}

// Logf makes use of log
func Logf(file string, line int, l level, format string, v ...interface{}) {
	fl := fmt.Sprintf("[%5s] %s:%d: ", levelString[l], file, line)
	lv := fmt.Sprintf(format, v...)
	// log.Print(fl, lv)
	str := fmt.Sprint(fl, lv)
	logChn <- str
}

// WithLevel logs with the level specified
func WithLevel(file string, line int, l level, v ...interface{}) {
	if l > logLevel {
		return
	}
	Log(file, line, l, v...)
}

// WithLevelf logs with the level specified
func WithLevelf(file string, line int, l level, format string, v ...interface{}) {
	if l > logLevel {
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
