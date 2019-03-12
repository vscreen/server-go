package player

import (
	"math"
	"time"
)

type timer struct {
	t        *time.Ticker
	position float64       // in seconds
	duration float64       // in seconds
	done     chan struct{} // tell ticker to stop
	c        chan<- float64
	C        <-chan float64
}

func newTimer(seconds int64) *timer {
	positionChan := make(chan float64)

	t := timer{
		t:        nil,
		position: 0,
		duration: float64(seconds),
		done:     nil,
		C:        positionChan,
		c:        positionChan,
	}

	t.play()
	return &t
}

// play does nothing if t is not nil. Else, play starts timer.
func (t *timer) play() {
	if t.t != nil {
		return
	}

	t.t = time.NewTicker(time.Second)
	t.done = make(chan struct{})
	go func(ticker *time.Ticker, done <-chan struct{}) {
	loop:
		for {
			select {
			case <-ticker.C:
			case <-done:
				break loop
			}

			t.position++
			curPos := t.position / t.duration // convert it to [0.0, 1.0]
			if curPos >= 1.0 {
				defer recover() // TODO! Avoid double close for c channel
				close(t.c)
				return
			}

			// If there's slow receiver, newest position will be disposed
			select {
			case t.c <- curPos:
			default:
			}
		}
	}(t.t, t.done)
}

// pause does nothing if t is nil. Else, pause set curDur with elapsed time
// from startTime, and set t to be nil
func (t *timer) pause() {
	if t.t == nil {
		return
	}

	t.t.Stop()
	close(t.done)
	t.t = nil
}

// seek sets curDur properly with pos is [0.0, 1.0].
// if t is not nil. seek will stop and reset timer with properly and
// set startTime to be current time.
func (t *timer) seek(pos float64) {
	running := false
	if t.t != nil {
		running = true
		t.t.Stop()
		close(t.done)
		t.t = nil
	}

	t.position = math.Round(pos * t.duration)

	// If the timer wasn't paused, resume the ticking
	if running {
		t.play()
	}
}

// stop does nothing if t is nil. Else, it'll stop the timer and no further
// methods should not be called.
func (t *timer) stop() {
	if t.t == nil {
		return
	}

	t.t.Stop()
	close(t.done)
	defer recover() // TODO! Avoid double close for c channel
	close(t.c)
}
