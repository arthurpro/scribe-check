package main

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Spinner renders a unicode spinner with elapsed time on a single line of an
// io.Writer (typically os.Stderr). Safe to Stop() multiple times.
type Spinner struct {
	w        io.Writer
	msg      string
	frames   []string
	interval time.Duration
	done     chan struct{}
	once     sync.Once
	wg       sync.WaitGroup
}

func NewSpinner(w io.Writer, msg string) *Spinner {
	return &Spinner{
		w:        w,
		msg:      msg,
		frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		interval: 100 * time.Millisecond,
		done:     make(chan struct{}),
	}
}

func (s *Spinner) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		start := time.Now()
		i := 0
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-s.done:
				return
			case <-ticker.C:
				elapsed := time.Since(start).Round(time.Second)
				fmt.Fprintf(s.w, "\r\x1b[2K  %s %s  %s  ", s.frames[i%len(s.frames)], s.msg, elapsed)
				i++
			}
		}
	}()
}

// Stop signals the render goroutine, waits for it to exit, then clears the
// spinner line. Safe to call multiple times.
func (s *Spinner) Stop() {
	s.once.Do(func() {
		close(s.done)
		s.wg.Wait()
		fmt.Fprint(s.w, "\r\x1b[2K")
	})
}
