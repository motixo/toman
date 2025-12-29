package spinner

import (
	"fmt"
	"sync"
	"time"
)

type Spinner struct {
	chars []rune
	delay time.Duration
	done  chan struct{}
	once  sync.Once
}

func New() *Spinner {
	return &Spinner{
		chars: []rune{'|', '/', '-', '\\'},
		delay: 100 * time.Millisecond,
		done:  make(chan struct{}),
	}
}

func (s *Spinner) Start() {
	go func() {
		i := 0
		for {
			select {
			case <-s.done:
				fmt.Print("\r\033[K")
				return
			default:
				fmt.Printf("\r%c Fetching data...", s.chars[i%len(s.chars)])
				time.Sleep(s.delay)
				i++
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.once.Do(func() {
		s.done <- struct{}{}
		close(s.done)
	})
}
