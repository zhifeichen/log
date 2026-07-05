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
	withBody   bool
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
// level(ascending order): trace, debug, info, warn, error, fatal
// default debug
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
// default ./log.log
func Filename(p string) Option {
	return func(o *Options) {
		o.filename = p
	}
}

// MaxSize of log(megabytes)
// default 10M
func MaxSize(s int) Option {
	return func(o *Options) {
		o.maxSize = s
	}
}

// MaxAge of log(days)
// default 30
func MaxAge(a int) Option {
	return func(o *Options) {
		o.maxAge = a
	}
}

// MaxBackups of log
// default 20
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

// WithBody set whether to log HTTP request/response bodies
func WithBody(b bool) Option {
	return func(o *Options) {
		o.withBody = b
	}
}

// toZapConfig converts Options to internal zap config
func (o Options) toZapConfig() zapConfig {
	levelStr := strings.ToLower(levelString[o.level])
	return zapConfig{
		level:      levelStr,
		filename:   o.filename,
		maxSize:    o.maxSize,
		maxAge:     o.maxAge,
		maxBackups: o.maxBackups,
		writers:    o.writers,
		withBody:   o.withBody,
	}
}

// zapConfig is the internal configuration for creating a zap-based Logger
type zapConfig struct {
	level      string
	filename   string
	maxSize    int
	maxAge     int
	maxBackups int
	writers    []io.Writer
	withBody   bool
}
