package operation

import (
	"errors"
	"sync"
)

var (
	ErrMutexChannelUnready = errors.New("channel is not prepared")
)

type AtomicMutex struct {
	fn   func()
	mu   sync.Mutex    // lock to avoid simultaneous calls
	done chan struct{} // to signal when the operation is done
}

func (am *AtomicMutex) Do() {
	am.mu.Lock()
	defer am.mu.Unlock()

	// prepare the channel
	am.done = make(chan struct{})
	defer close(am.done)

	am.fn()
}

func (am *AtomicMutex) Done() (<-chan struct{}, error) {
	if am.done == nil {
		return nil, ErrMutexChannelUnready
	}

	return am.done, nil
}
