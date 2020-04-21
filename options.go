package log

import (
	"io"
	"strings"
)

// Options log options
type Options struct {
	level      level
	filename   string
	maxSize    int
	maxAge     int
	maxBackups int
	writers    []io.Writer
}

// Option log option
type Option func(*Options)

// NewOptions new log options
func NewOptions(options ...Option) Options {
	opts := Options{
		level:      LevelDebug,
		filename:   "./log.log",
		maxSize:    10, // default 10M
		maxAge:     30, // default 30 days
		maxBackups: 20, // default 20 backups
	}
	for _, o := range options {
		o(&opts)
	}
	return opts
}

// Level of log
func Level(l string) Option {
	newl := LevelDebug
	switch strings.ToLower(l) {
	case "fatal":
		newl = LevelFatal
	case "error":
		newl = LevelError
	case "warn":
		newl = LevelWarn
	case "info":
		newl = LevelInfo
	case "debug":
		newl = LevelDebug
	case "trace":
		newl = LevelTrace
	}
	return func(o *Options) {
		o.level = newl
	}
}

// Filename set log file name
func Filename(p string) Option {
	return func(o *Options) {
		o.filename = p
	}
}

// MaxSize of log
func MaxSize(s int) Option {
	return func(o *Options) {
		o.maxSize = s
	}
}

// MaxAge of log
func MaxAge(a int) Option {
	return func(o *Options) {
		o.maxAge = a
	}
}

// MaxBackups of log
func MaxBackups(b int) Option {
	return func(o *Options) {
		o.maxBackups = b
	}
}

// Writers set multi writer
func Writers(w ...io.Writer) Option {
	return func(o *Options) {
		if o.writers == nil {
			o.writers = make([]io.Writer, 0)
		}
		o.writers = append(o.writers, w...)
	}
}
