package operation

import (
	"reflect"
	"testing"
	"time"
)

func TestNewAtomicState(t *testing.T) {
	t.Run("success to create an instance of atomic state", func(t *testing.T) {
		fn := func() {
			// do something
		}

		as, err := NewAtomicState(fn)

		if err != nil {
			t.Errorf("expected no error, got %s", err.Error())
			return
		}
		expectedPtrFn, ptrFn := reflect.ValueOf(fn).Pointer(), reflect.ValueOf(as.fn).Pointer()
		if expectedPtrFn != ptrFn {
			t.Errorf("expected fn with pointer %d, got %d", expectedPtrFn, ptrFn)
			return
		}
		if as.state != 0 {
			t.Errorf("expected state to be initialized as 0, got %d", as.state)
			return
		}
		if as.done != nil {
			t.Error("expected channel done to be nil, got initialized channel")
			return
		}
	})

	t.Run("failed to create an instance of atomic state - fn is nil", func(t *testing.T) {
		fn := (func())(nil)

		as, err := NewAtomicState(fn)

		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		if err != ErrStateNewFunctionNil {
			t.Errorf("expected error %s, got %s", ErrStateNewFunctionNil.Error(), err.Error())
			return
		}
		if as != nil {
			t.Errorf("expected atomic state to be nil, got %v", as)
			return
		}
	})
}

func TestAtomicState_Do(t *testing.T) {
	t.Run("success to run operation, in concurrent ambience and signaling", func(t *testing.T) {
		am := &AtomicState{
			fn: func() {
				time.Sleep(time.Second)
			},
			state: 0,
			done:  nil,
		}

		done := make(chan error, 1)
		go func() {
			defer close(done)
			err := am.Do()
			if err != nil {
				done <- err
				return
			}
		}()

		time.Sleep(time.Millisecond * 100)

		for i := 0; i < 2; i++ {
			err := am.Do()

			if err == nil {
				t.Error("expected error, got nil")
				return
			}
			if err != ErrStateDoUnavailable {
				t.Errorf("expected error %s, got %s", ErrStateDoUnavailable.Error(), err.Error())
				return
			}
		}

		select {
		case <-time.After(time.Second):
			t.Error("expected function to be done running, but did not")
			return
		case err, ok := <-done:
			if ok {
				t.Errorf("expected no error, got %s", err.Error())
				return
			}
		}

		select {
		case <-time.After(time.Millisecond * 250):
			t.Error("expected channel done to be closed, but did not")
			return
		case <-am.done:
			// success
		}
	})
}

func TestAtomicState_DoAsync(t *testing.T) {
	t.Run("success to run operation, in concurrent ambience and signaling", func(t *testing.T) {
		am := &AtomicState{
			fn: func() {
				time.Sleep(time.Second)
			},
			state: 0,
			done:  nil,
		}

		err := am.DoAsync()
		if err != nil {
			t.Errorf("expected no error, got %s", err.Error())
			return
		}

		for i := 0; i < 2; i++ {
			err := am.DoAsync()

			if err == nil {
				t.Error("expected error, got nil")
				return
			}
			if err != ErrStateDoUnavailable {
				t.Errorf("expected error %s, got %s", ErrStateDoUnavailable.Error(), err.Error())
				return
			}
		}

		// wait at least 250 ms until the go routine beggins
		// - otherwise the am.done is considered nil and will only read when time.After is done, even the close in the go routine after was initialized and then closed properly.
		time.Sleep(time.Millisecond * 250)
		select {
		case <-time.After(time.Second):
			t.Error("expected channel done to be closed, but did not")
			return
		case <-am.done:
			// success
		}
	})
}

func TestAtomicState_Done(t *testing.T) {
	t.Run("success to share an only-read channel, done is prepare", func(t *testing.T) {
		am := &AtomicState{
			fn:    nil,
			state: 0,
			done:  make(chan struct{}),
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
		am := &AtomicState{
			fn:    nil,
			state: 0,
			done:  nil,
		}

		ch, err := am.Done()

		if err == nil {
			t.Error("expected error, got nil")
			return
		}
		if err != ErrStateChannelUnready {
			t.Errorf("expected err %s, got %s", ErrStateChannelUnready, err.Error())
			return
		}
		if ch != nil {
			t.Errorf("expected channel to be nil, got %v", ch)
			return
		}
	})
}
