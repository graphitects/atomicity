package operation

import (
	"testing"
	"time"
)

func TestAtomicState_Do_Done(t *testing.T) {
	t.Run("success to execute Do() and receive the signal from Done()", func(t *testing.T) {
		as := &AtomicState{
			fn:    func() { time.Sleep(time.Second) },
			state: 0,
			done:  nil,
		}

		done := make(chan error, 1)
		go func() {
			defer close(done)

			err := as.Do()
			if err != nil {
				done <- err
				return
			}
		}()

		time.Sleep(time.Millisecond * 100)

		subs := [2]<-chan struct{}{nil, nil}
		for i := 0; i < 2; i++ {
			var err error
			subs[i], err = as.Done()
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
		case err, ok := <-done:
			if ok {
				t.Errorf("expected channel to be closed, but got error msg %s", err.Error())
				return
			}
		}
	})
}
