package player

import (
	"testing"
	"time"
)

// allows only 2 ms off
const eps = 1000000 // 1000000ns = 1000μs = 1ms

func TestPlay(t *testing.T) {
	start := time.Now()
	done := make(chan struct{})
	expected := time.Second.Nanoseconds()

	timer := newTimer(time.Second*1, func() {
		done <- struct{}{}
	})

	timer.play()
	<-done

	elapsed := time.Since(start).Nanoseconds()
	if !(expected-eps <= elapsed && elapsed <= expected+eps) {
		t.Error("play failed")
	}
}

func TestSeek(t *testing.T) {
	start := time.Now()
	done := make(chan struct{})
	expected := time.Second.Nanoseconds()

	timer := newTimer(time.Second*10, func() {
		done <- struct{}{}
	})

	timer.seek(0.9)
	timer.play()
	<-done

	elapsed := time.Since(start).Nanoseconds()
	if !(expected-eps <= elapsed && elapsed <= expected+eps) {
		t.Error("seek failed")
	}
}

func TestPause(t *testing.T) {
	start := time.Now()
	done := make(chan struct{})
	expected := (time.Second * 3).Nanoseconds()

	timer := newTimer(time.Second*2, func() {
		done <- struct{}{}
	})

	timer.play()
	<-time.After(time.Second)
	timer.pause()
	<-time.After(time.Second)
	timer.play()
	<-done

	elapsed := time.Since(start).Nanoseconds()
	if !(expected-eps <= elapsed && elapsed <= expected+eps) {
		t.Error("pause failed")
	}
}
