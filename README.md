# Atomicity Package

The `atomicity` package provides concurrency-safe utilities for ensuring controlled execution of critical sections. It offers mechanisms for both synchronous and asynchronous operations, with built-in signaling to notify when operations are complete.

## Features

- Ensures only one operation runs at a time using atomic state control.
- Provides synchronous (`Do`) and asynchronous (`DoAsync`) execution methods.
- Includes a signaling mechanism (`Done`) to notify when an operation has completed.
- Built-in error handling for invalid states and uninitialized channels.

## Installation

```bash
go get github.com/graphitects/atomicity
```

## Usage

### AtomicState

The `AtomicState` structure allows you to control and monitor operations with concurrency safety.

### API Overview

#### `AtomicState` Structure

```go
// AtomicState is a concurrency-safe structure that ensures an operation can only run one at a time.
// It provides synchronous and asynchronous execution methods and a signaling mechanism
// to notify when the operation is complete.

type AtomicState struct {
    fn    func()        // The function to be executed safely.
    state uint32        // Atomic state used to control access to the operation.
    done  chan struct{} // Channel used to signal when the operation is complete.
}
```

#### `Do` Method

Executes the function synchronously, ensuring that only one operation can run at a time.

```go
// Do executes the function synchronously, ensuring that only one operation
// can run at a time. It uses atomic state control to determine availability,
// reinitializes the `done` channel with each invocation, and signals completion
// by closing the channel.
//
// Returns:
// - `nil` if the operation is successfully executed.
// - An error if another operation is already in progress.
func (am *AtomicState) Do() error
```

#### `DoAsync` Method

Executes the function asynchronously, ensuring that only one operation can run at a time.

```go
// DoAsync executes the function asynchronously, ensuring that only one operation
// can run at a time. It uses atomic state control to determine availability,
// reinitializes the `done` channel with each invocation, and signals completion
// by closing the channel.
//
// Note:
// - The `done` channel is created within the goroutine. Consumers should
//   avoid accessing it until the goroutine has started to prevent race conditions.
//
// Returns:
// - `nil` if the operation is successfully scheduled.
// - An error if another operation is already in progress.
func (am *AtomicState) DoAsync() error
```

#### `Done` Method

Provides access to the `done` channel, signaling when the operation has completed.

```go
// Done provides access to the `done` channel, which signals when the operation
// has completed. It leverages Go's native behavior of broadcasting to all listeners
// when a channel is closed. If the `done` channel is uninitialized, an error is returned.
//
// Note:
// - For operations executed with `DoAsync`, ensure that the `done` channel is accessed
//   only after the asynchronous operation has started to avoid receiving an uninitialized channel.
//
// Returns:
// - The `done` channel (read-only) for signaling operation completion.
// - An error if the `done` channel has not been prepared.
func (am *AtomicState) Done() (<-chan struct{}, error)
```

### Example Usage

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/graphitects/atomicity/operation"
)

func main() {
    // Define a function to be executed.
    myFunc := func() {
        fmt.Println("Operation started")
        time.Sleep(2 * time.Second)
        fmt.Println("Operation completed")
    }

    // Create an instance of AtomicState using the `NewAtomicState` constructor.
    am, err := operation.NewAtomicState(myFunc)
    if err != nil {
        log.Fatal(err)
    }

    // Run the operation synchronously.
    if err := am.Do(); err != nil {
        log.Fatal(err)
    }

    // Wait for completion.
    done, err := am.Done()
    if err != nil {
        log.Fatal(err)
    }
    <-done

    fmt.Println("Synchronous operation finished")

    // Run the operation asynchronously.
    if err := am.DoAsync(); err != nil {
        log.Fatal(err)
    }

    // Wait for completion.
    done, err = am.Done()
    if err != nil {
        log.Fatal(err)
    }
    <-done

    fmt.Println("Asynchronous operation finished")
}
```

## Errors

### `ErrStateDoUnavailable`
Returned when an operation is already in progress.

### `ErrStateChannelUnready`
Returned when the `done` channel is not initialized.

## Contributing

Feel free to contribute by opening issues or submitting pull requests. Ensure all contributions are well-tested and documented.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
