package player

import (
	"strconv"
	"testing"
	"time"
)

// allows only 1 ms off
const eps = 1000000 // 1000000ns = 1000μs = 1ms
const cases = 10

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
				t.Errorf("expected: %d, but got: %d", expected, elapsed)
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
				t.Errorf("expected: %d, but got: %d", expected, elapsed)
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
				t.Errorf("expected: %d, but got: %d", expected, elapsed)
			}
		})
	}
}

func BenchmarkPause(b *testing.B) {

	for n := 0; n < b.N; n++ {
		done := make(chan struct{})
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
	}
}
