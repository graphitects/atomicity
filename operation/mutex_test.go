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
