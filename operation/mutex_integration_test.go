package operation

import (
	"sync"
	"testing"
	"time"
)

func TestAtomicMutex_Do_Done(t *testing.T) {
	t.Run("success to Do the function and receive the signal from done", func(t *testing.T) {
		am := &AtomicMutex{
			fn:   func() { time.Sleep(time.Second) },
			mu:   sync.Mutex{},
			done: nil,
		}

		// execute do function
		done := make(chan struct{})
		go func() {
			defer close(done)
			am.Do()
		}()

		time.Sleep(time.Millisecond * 100)

		subs := [2]<-chan struct{}{nil, nil}
		for i := 0; i < 2; i++ {
			var err error
			subs[i], err = am.Done()
			if err != nil {
				t.Errorf("expected no error, got %s", err.Error())
				return
			}
		}

		for i := 0; i < 2; i++ {
			select {
			case <-time.After(time.Second):
				t.Error("expected channel to be closed, but did not")
				return
			case <-subs[i]:
				// success
			}
		}

		select {
		case <-time.After(time.Millisecond * 100):
			t.Error("expected main call to be done, but did not")
			return
		case <-done:
			return
		}
	})
}
