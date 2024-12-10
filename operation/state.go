package operation

import (
	"errors"
	"sync/atomic"
)

var (
	ErrStateDoUnavailable  = errors.New("operation is not available")
	ErrStateChannelUnready = errors.New("channel is not prepared")
)

type AtomicState struct {
	fn    func()
	state uint32
	done  chan struct{} // to signal when the operation is done
}

func (am *AtomicState) Do() error {
	if !atomic.CompareAndSwapUint32(&am.state, 0, 1) {
		return ErrStateDoUnavailable
	}
	defer atomic.StoreUint32(&am.state, 0)

	// prepare the channel
	if am.done == nil {
		am.done = make(chan struct{})
	}
	defer close(am.done)

	am.fn()
	return nil
}

func (am *AtomicState) DoAsync() error {
	if !atomic.CompareAndSwapUint32(&am.state, 0, 1) {
		return ErrStateDoUnavailable
	}

	go func() {
		defer atomic.StoreUint32(&am.state, 0)
		defer close(am.done)

		// prepare channel
		if am.done == nil {
			am.done = make(chan struct{})
		}
		am.fn()
	}()

	return nil
}

func (am *AtomicState) Done() (<-chan struct{}, error) {
	if am.done == nil {
		return nil, ErrStateChannelUnready
	}

	return am.done, nil
}
