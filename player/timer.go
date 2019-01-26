package player

import (
	"sync"
	"time"
)

type finishCallback func()

type timer struct {
	t         *time.Timer
	origDur   time.Duration
	curDur    time.Duration
	startTime time.Time
	f         finishCallback
	m         *sync.Mutex
}

func newTimer(d time.Duration, f finishCallback) *timer {
	return &timer{
		t:       nil,
		origDur: d,
		curDur:  d,
		f:       f,
		m:       &sync.Mutex{},
	}
}

// play does nothing if t is not nil. Else, play starts timer
// and set startTime to be the current time
func (t *timer) play() {
	t.m.Lock()
	defer t.m.Unlock()
	if t.t != nil {
		return
	}

	t.t = time.AfterFunc(t.curDur, t.f)
	t.startTime = time.Now()
}

// pause does nothing if t is nil. Else, pause set curDur with elapsed time
// from startTime, and set t to be nil
func (t *timer) pause() {
	t.m.Lock()
	defer t.m.Unlock()
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
	t.m.Lock()
	defer t.m.Unlock()
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
	t.m.Lock()
	defer t.m.Unlock()
	if t.t == nil {
		return
	}

	t.t.Stop()
	t.t = nil
}
