package operation

import (
	"sync"
	"testing"
	"time"
)

func TestAtomicMutex_Do(t *testing.T) {
	t.Run("success to run operation, in concurrent ambience and signaling", func(t *testing.T) {
		am := &AtomicMutex{
			fn: func() {
				time.Sleep(time.Second)
			},
			mu:   sync.Mutex{},
			done: make(chan struct{}),
		}

		// main call
		done := make(chan struct{})
		go func() {
			defer close(done)
			am.Do()
		}()

		time.Sleep(time.Millisecond * 100)

		doneSecondCall := make(chan struct{})
		go func() {
			defer close(doneSecondCall)
			am.Do()
		}()
		doneThirdCall := make(chan struct{})
		go func() {
			defer close(doneThirdCall)
			am.Do()
		}()

		select {
		case <-time.After(time.Millisecond * 500):
			// success
		case <-doneSecondCall:
			t.Error("expected second call to get locked, but did not")
			return
		case <-doneThirdCall:
			t.Error("expected third call to get locked, but did not")
			return
		}

		select {
		case <-time.After(time.Millisecond * 1500):
			t.Error("expected main call to be done, but did not")
			return
		case <-done:
			return
		}
	})
}

func TestAtomicMutex_Done(t *testing.T) {
	t.Run("success to share an only-read channel, done is prepare", func(t *testing.T) {
		am := &AtomicMutex{
			fn:   nil,
			mu:   sync.Mutex{},
			done: make(chan struct{}),
		}

		ch, err := am.Done()

		if err != nil {
			t.Errorf("expected no error, got %s", err.Error())
			return
		}
		expected := (<-chan struct{})(am.done)
		if expected != ch {
			t.Errorf("expected channel %v, got %v", expected, ch)
			return
		}
	})

	t.Run("failure to share an only-read channel, done is not prepared", func(t *testing.T) {
		am := &AtomicMutex{
			fn:   nil,
			mu:   sync.Mutex{},
			done: nil,
		}

		ch, err := am.Done()

		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		if err != ErrMutexChannelUnready {
			t.Errorf("expected err %s, got %s", ErrMutexChannelUnready, err.Error())
			return
		}
		if ch != nil {
			t.Errorf("expected channel to be nil, got %v", ch)
			return
		}
	})
}
