package log

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger the logger struct
type Logger struct {
	logger    *log.Logger
	logLevel  level
	oldWriter io.Writer
	bufWriter *bufio.Writer
	debounce  debounce
	chn       chan string
}

// New init log
func New(o Options) *Logger {
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

	bufWriter := bufio.NewWriterSize(w, 64*1024)
	l := log.New(bufWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	ll := &Logger{
		logger:    l,
		logLevel:  o.level,
		bufWriter: bufWriter,
		chn:       make(chan string, 50),
	}
	ll.debounce = newDebouncer(time.Second, func() {
		ll.bufWriter.Flush()
	})
	go ll.loop()
	return ll
}

func (logger *Logger) Discard() {
	logger.oldWriter = logger.logger.Writer()
	logger.logger.SetOutput(ioutil.Discard)
}

func (logger *Logger) ResumeWriter() {
	if logger.oldWriter != nil {
		logger.logger.SetOutput(logger.oldWriter)
		logger.oldWriter = nil
	}
}

func (logger *Logger) Flush() error {
	if logger.bufWriter != nil {
		return logger.bufWriter.Flush()
	}
	return nil
}

func (logger *Logger) flush() {
	logger.debounce()
}

func (logger *Logger) loop() {
	for {
		logger.flush()
		l, ok := <-logger.chn
		if !ok {
			return
		}
		if logger.bufWriter != nil && len(l) > logger.bufWriter.Available() {
			logger.bufWriter.Flush()
		}
		logger.logger.Print(l)
	}
}

// Log makes use of log
func (logger *Logger) Log(file string, line int, l level, v ...interface{}) {
	fl := fmt.Sprintf("[%5s] %s:%d:", levelString[l], file, line)
	lv := fmt.Sprint(v...)
	// size := len(fl) + len(lv)
	// if logger.bufWriter != nil && size > logger.bufWriter.Available() {
	// 	logger.Flush()
	// }
	// logger.logger.Println(fl, lv)
	// logger.flush()
	buf := fmt.Sprintln(fl, lv)
	logger.chn <- buf
}

// Logf makes use of log
func (logger *Logger) Logf(file string, line int, l level, format string, v ...interface{}) {
	fl := fmt.Sprintf("[%5s] %s:%d: ", levelString[l], file, line)
	lv := fmt.Sprintf(format, v...)
	// size := len(fl) + len(lv)
	// if logger.bufWriter != nil && size > logger.bufWriter.Available() {
	// 	logger.Flush()
	// }
	// logger.logger.Print(fl, lv)
	// logger.flush()
	buf := fmt.Sprint(fl, lv)
	logger.chn <- buf
}

// WithLevel logs with the level specified
func (logger *Logger) WithLevel(file string, line int, l level, v ...interface{}) {
	if l > logger.logLevel {
		return
	}
	logger.Log(file, line, l, v...)
}

// WithLevelf logs with the level specified
func (logger *Logger) WithLevelf(file string, line int, l level, format string, v ...interface{}) {
	if l > logger.logLevel {
		return
	}
	logger.Logf(file, line, l, format, v...)
}

// Trace provides trace level logging
func (logger *Logger) Trace(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevel(filepath.Base(file), line, LevelTrace, v...)
}

// Tracef provides trace level logging
func (logger *Logger) Tracef(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevelf(filepath.Base(file), line, LevelTrace, format, v...)
}

// Debug provides debug level logging
func (logger *Logger) Debug(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevel(filepath.Base(file), line, LevelDebug, v...)
}

// Debugf provides debug level logging
func (logger *Logger) Debugf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevelf(filepath.Base(file), line, LevelDebug, format, v...)
}

// Info provides info level logging
func (logger *Logger) Info(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevel(filepath.Base(file), line, LevelInfo, v...)
}

// Infof provides info level logging
func (logger *Logger) Infof(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevelf(filepath.Base(file), line, LevelInfo, format, v...)
}

// Warn provides warn level logging
func (logger *Logger) Warn(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevel(filepath.Base(file), line, LevelWarn, v...)
}

// Warnf provides warn level logging
func (logger *Logger) Warnf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevelf(filepath.Base(file), line, LevelWarn, format, v...)
}

// Error provides error level logging
func (logger *Logger) Error(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevel(filepath.Base(file), line, LevelError, v...)
}

// Errorf provides error level logging
func (logger *Logger) Errorf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevelf(filepath.Base(file), line, LevelError, format, v...)
}

// Fatal provides fatal level logging
func (logger *Logger) Fatal(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevel(filepath.Base(file), line, LevelFatal, v...)
	os.Exit(1)
}

// Fatalf provides fatal level logging
func (logger *Logger) Fatalf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logger.WithLevelf(filepath.Base(file), line, LevelFatal, format, v...)
	os.Exit(1)
}

// Writer get log writer
func (logger *Logger) Writer() io.Writer {
	return logger.logger.Writer()
}
