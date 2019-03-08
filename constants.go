package comshim

import "errors"

var (
	// ErrNegativeCounter is returned when the internal counter of a shim drops
	// below zero. This may indicate that Done() has been called more than once
	// for the same object.
	ErrNegativeCounter = errors.New("component object model shim counter has dropped below zero")

	// ErrAlreadyInitialized is returned when a shim finds itself on a thread
	// that has already been initialized. This probably indicates that some
	// previous goroutine failed to lock the OS thread or failed to call
	// CoUninitialize when it should have.
	ErrAlreadyInitialized = errors.New("component object model shim thread has already been initialized")
)
