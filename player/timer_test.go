package player

import (
	"strconv"
	"testing"
	"time"
)

// allows only 2 ms off
const eps = 1000000 // 1000000ns = 1000μs = 1ms
const cases = 100

func TestPlay(t *testing.T) {
	for i := 0; i < cases; i++ {
		i := i
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			start := time.Now()
			done := make(chan struct{})
			expected := time.Second.Nanoseconds()

			timer := newTimer(1)
			go func() {
				for range timer.C {
				}
				done <- struct{}{}
			}()

			timer.play()
			<-done

			elapsed := time.Since(start).Nanoseconds()
			if !(expected-eps <= elapsed && elapsed <= expected+eps) {
				t.Error("play failed")
			}
		})
	}
}

func TestSeek(t *testing.T) {
	for i := 0; i < cases; i++ {
		i := i
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			start := time.Now()
			done := make(chan struct{})
			expected := (2 * time.Second).Nanoseconds()

			timer := newTimer(511)
			go func() {
				for range timer.C {
				}
				done <- struct{}{}
			}()

			timer.seek(float64(509) / 511)
			<-done

			elapsed := time.Since(start).Nanoseconds()
			if !(expected-eps <= elapsed && elapsed <= expected+eps) {
				t.Error("seek failed")
			}
		})
	}
}

func TestPause(t *testing.T) {
	for i := 0; i < cases; i++ {
		i := i
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			start := time.Now()
			done := make(chan struct{})
			expected := (time.Second * 3).Nanoseconds()

			timer := newTimer(2)
			go func() {
				for range timer.C {
				}
				done <- struct{}{}
			}()

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
		})
	}
}
