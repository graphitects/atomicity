package operation

import (
	"errors"
	"sync"
)

// Predefined error messages for common failure scenarios.
var (
	ErrNewFunctionNil      = errors.New("function can not be nil")
	ErrMutexChannelUnready = errors.New("channel is not prepared") // Error returned when the done channel is uninitialized.
)

// AtomicMutex is a concurrency-safe structure that ensures an operation
// can only run one at a time, and it provides a signaling mechanism
// to notify when the operation is complete.
type AtomicMutex struct {
	fn   func()        // The function to be executed safely.
	mu   sync.Mutex    // Mutex to prevent simultaneous calls to the Do method.
	done chan struct{} // Channel used to signal when the operation is complete.
}

// NewAtomicMutex creates a new instance of the `AtomicMutex` struct, which
// ensures that the function is executed atomically and provides a signaling
// mechanism to indicate when the operation is complete.
//
// Parameters:
// - fn: The function to be executed atomically.
//
// Returns:
// - A new instance of the `AtomicMutex` struct.
// - An error if the function is nil.
func NewAtomicMutex(fn func()) (*AtomicMutex, error) {
	if fn == nil {
		return nil, ErrNewFunctionNil
	}

	return &AtomicMutex{
		fn:   fn,
		mu:   sync.Mutex{},
		done: nil,
	}, nil
}

// Do locks the mutex to ensure that the function is executed safely.
// It reinitializes the `done` channel with each invocation, executes
// the function, and signals completion by closing the `done` channel.
func (am *AtomicMutex) Do() {
	am.mu.Lock()         // Lock the mutex to ensure exclusive access.
	defer am.mu.Unlock() // Unlock the mutex after the operation completes.

	am.done = make(chan struct{}) // Reinitialize the `done` channel for a new operation.
	defer close(am.done)          // Close the channel to broadcast completion to listeners.

	am.fn() // Execute the function.
}

// Done provides access to the `done` channel, which signals when the operation
// has completed. If `Do` has not been called or the `done` channel is uninitialized,
// an error is returned.
//
// Returns:
// - The `done` channel (read-only) for signaling operation completion.
// - An error if the `done` channel has not been prepared.
func (am *AtomicMutex) Done() (<-chan struct{}, error) {
	if am.done == nil {
		return nil, ErrMutexChannelUnready // Return an error if the channel is not initialized.
	}

	return am.done, nil // Return the `done` channel for consumers to listen for completion.
}
