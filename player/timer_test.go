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

	timer := newTimer(1, func() {
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
	expected := (2 * time.Second).Nanoseconds()

	timer := newTimer(511, func() {
		done <- struct{}{}
	})

	timer.seek(float64(509) / 511)
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

	timer := newTimer(2, func() {
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
