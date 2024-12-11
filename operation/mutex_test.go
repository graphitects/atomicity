package operation

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewAtomicMutex(t *testing.T) {
	t.Run("success to create an instance of atomic mutex", func(t *testing.T) {
		fn := func() {
			// do something
		}

		am, err := NewAtomicMutex(fn)

		if err != nil {
			t.Errorf("expected no error, got %s", err.Error())
			return
		}
		expectedPtrFn, ptrFn := reflect.ValueOf(fn).Pointer(), reflect.ValueOf(am.fn).Pointer()
		if expectedPtrFn != ptrFn {
			t.Errorf("expected fn with pointer %d, got %d", expectedPtrFn, ptrFn)
			return
		}
		if am.done != nil {
			t.Error("expected channel done to be nil, got initialized channel")
			return
		}
	})

	t.Run("failed to create an instance of atomic mutex - fn is nil", func(t *testing.T) {
		fn := (func())(nil)

		am, err := NewAtomicMutex(fn)

		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		if err != ErrNewFunctionNil {
			t.Errorf("expected error %s, got %s", ErrNewFunctionNil.Error(), err.Error())
			return
		}
		if am != nil {
			t.Errorf("expected atomic mutex to be nil, got %v", am)
			return
		}
	})
}

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
