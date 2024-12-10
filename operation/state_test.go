package operation

import (
	"testing"
	"time"
)

func TestAtomicState_Do(t *testing.T) {
	t.Run("success to run operation, in concurrent ambience and signaling", func(t *testing.T) {
		am := &AtomicState{
			fn: func() {
				time.Sleep(time.Second)
			},
			state: 0,
			done:  make(chan struct{}),
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
	})
}
