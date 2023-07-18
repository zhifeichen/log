package log

import (
	"time"
)

type debounce func()

type debouncer struct {
	after time.Duration
	timer *time.Timer
}

func (d *debouncer) do() {
	d.timer.Reset(d.after)
}

func newDebouncer(after time.Duration, f func()) debounce {
	d := &debouncer{after: after, timer: time.AfterFunc(after, f)}
	return func() {
		d.do()
	}
}
