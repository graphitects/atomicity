package operation

import (
	"errors"
	"sync"
)

// Predefined error messages for common failure scenarios.
var (
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

// Do locks the mutex to ensure that the function is executed safely.
// It initializes the `done` channel, executes the function, and then
// signals completion by closing the `done` channel.
func (am *AtomicMutex) Do() {
	am.mu.Lock()         // Lock the mutex to ensure exclusive access.
	defer am.mu.Unlock() // Unlock the mutex after the operation completes.

	// Initialize the `done` channel to signal the completion of the operation.
	if am.done == nil {
		am.done = make(chan struct{})
	}
	defer close(am.done) // Close the channel to broadcast completion to listeners.

	am.fn() // Execute the function.
}

// Done provides access to the `done` channel, which signals when the operation
// has completed. It leverages Go's native behavior of broadcasting to all listeners
// when a channel is closed. If the `done` channel is uninitialized, an error is returned.
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
