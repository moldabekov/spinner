package spinner

import (
	"fmt"
	"sync/atomic"
	"time"
)

const ClearLine = "\r\033[K"

var (
	Box1    = `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`
	Box2    = `⠄⠆⠇⠋⠙⠸⠰⠠⠰⠸⠙⠋⠇⠆`
	Spin1   = `◴◷◶◵`
	Spin2   = `◰◳◲◱`
	Spin3   = `←↑→↓`
	Default = Box1
)

type Spinner struct {
	frames []rune
	pos    int
	active uint64
	text   string
}

func New(text string) *Spinner {
	s := &Spinner{
		text: ClearLine + text,
	}
	s.Set(Default)
	return s
}

// Set frames to the given string which must not use spaces.
func (s *Spinner) Set(frames string) {
	s.frames = []rune(frames)
}

// Show the spinner
func (s *Spinner) Start() *Spinner {
	if atomic.LoadUint64(&s.active) > 0 {
		return s
	}
	atomic.StoreUint64(&s.active, 1)
	go func() {
		for atomic.LoadUint64(&s.active) > 0 {
			fmt.Printf(s.text, s.next())
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return s
}

// Hide the spinner
func (s *Spinner) Stop() bool {
	if x := atomic.SwapUint64(&s.active, 0); x > 0 {
		fmt.Printf(ClearLine)
		return true
	}
	return false
}

func (s *Spinner) next() string {
	r := s.frames[s.pos%len(s.frames)]
	s.pos++
	return string(r)
}
