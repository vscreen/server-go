package player

import (
	"time"
)

type finishCallback func()

type timer struct {
	t         *time.Timer
	origDur   time.Duration
	curDur    time.Duration
	startTime time.Time
	f         finishCallback
}

func newTimer(seconds int64, f finishCallback) *timer {
	d := time.Second * time.Duration(seconds)

	t := timer{
		t:       nil,
		origDur: d,
		curDur:  d,
		f:       f,
	}

	t.play()
	return &t
}

// play does nothing if t is not nil. Else, play starts timer
// and set startTime to be the current time
func (t *timer) play() {
	if t.t != nil {
		return
	}

	t.t = time.AfterFunc(t.curDur, t.f)
	t.startTime = time.Now()
}

// pause does nothing if t is nil. Else, pause set curDur with elapsed time
// from startTime, and set t to be nil
func (t *timer) pause() {
	if t.t == nil {
		return
	}

	t.t.Stop()
	elapsed := time.Since(t.startTime)
	t.curDur = t.curDur - elapsed
	t.t = nil
}

// seek sets curDur properly with pos is [0.0, 1.0].
// if t is not nil. seek will stop and reset timer with properly and
// set startTime to be current time.
func (t *timer) seek(pos float64) {
	t.curDur = t.origDur - time.Duration(pos*float64(t.origDur))

	if t.t != nil {
		t.t.Stop()
		t.t.Reset(t.curDur)
		t.startTime = time.Now()
	}
}

// stop does nothing if t is nil. Else, it'll stop the timer to avoid
// callback leak and set t to be nil
func (t *timer) stop() {
	if t.t == nil {
		return
	}

	t.t.Stop()
	t.t = nil
}

// pos gets current position in video in range of [0.0, 1.0]
func (t *timer) pos() float64 {
	return (t.origDur - t.curDur).Seconds() / t.origDur.Seconds()
}
