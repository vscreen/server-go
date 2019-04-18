package player

import (
	"testing"
	"time"
)

// allows only 5 ms off
const eps = 5000000

func TestPlay(t *testing.T) {
	start := time.Now()
	done := make(chan struct{})
	expected := time.Second.Nanoseconds()

	timer := newTimer(func() {
		done <- struct{}{}
	})

	timer.new(1)
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

	timer := newTimer(func() {
		done <- struct{}{}
	})

	timer.new(511)
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

	timer := newTimer(func() {
		done <- struct{}{}
	})

	timer.new(2)
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
