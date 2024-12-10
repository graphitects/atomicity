package operation

import (
	"errors"
	"sync/atomic"
)

// Predefined error messages for common failure scenarios.
var (
	ErrStateDoUnavailable  = errors.New("operation is not available") // Error returned when the operation is already in progress.
	ErrStateChannelUnready = errors.New("channel is not prepared")    // Error returned when the done channel is uninitialized.
)

// AtomicState is a concurrency-safe structure that ensures an operation can only run one at a time.
// It provides synchronous and asynchronous execution methods and a signaling mechanism
// to notify when the operation is complete.
type AtomicState struct {
	fn    func()        // The function to be executed safely.
	state uint32        // Atomic state used to control access to the operation.
	done  chan struct{} // Channel used to signal when the operation is complete.
}

// Do executes the function synchronously, ensuring that only one operation
// can run at a time. It uses atomic state control to determine availability,
// initializes the `done` channel, and signals completion by closing the channel.
func (am *AtomicState) Do() error {
	if !atomic.CompareAndSwapUint32(&am.state, 0, 1) {
		return ErrStateDoUnavailable // Return an error if the operation is already in progress.
	}
	defer atomic.StoreUint32(&am.state, 0) // Reset the state to allow future executions.

	am.done = make(chan struct{})
	defer close(am.done) // Close the channel to broadcast completion to listeners.

	am.fn() // Execute the function.
	return nil
}

// DoAsync executes the function asynchronously, ensuring that only one operation
// can run at a time. It uses atomic state control to determine availability,
// initializes the `done` channel, and signals completion by closing the channel.
func (am *AtomicState) DoAsync() error {
	if !atomic.CompareAndSwapUint32(&am.state, 0, 1) {
		return ErrStateDoUnavailable // Return an error if the operation is already in progress.
	}

	go func() {
		defer atomic.StoreUint32(&am.state, 0) // Reset the state to allow future executions.
		defer close(am.done)                   // Close the channel to broadcast completion to listeners.

		am.done = make(chan struct{})
		am.fn() // Execute the function asynchronously.
	}()

	return nil
}

// Done provides access to the `done` channel, which signals when the operation
// has completed. It leverages Go's native behavior of broadcasting to all listeners
// when a channel is closed. If the `done` channel is uninitialized, an error is returned.
//
// Returns:
// - The `done` channel (read-only) for signaling operation completion.
// - An error if the `done` channel has not been prepared.
func (am *AtomicState) Done() (<-chan struct{}, error) {
	if am.done == nil {
		return nil, ErrStateChannelUnready // Return an error if the channel is not initialized.
	}

	return am.done, nil // Return the `done` channel for consumers to listen for completion.
}
