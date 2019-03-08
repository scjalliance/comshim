comshim [![GoDoc](https://godoc.org/github.com/scjalliance/comshim?status.svg)](https://godoc.org/github.com/scjalliance/comshim)
====

The comshim package provides a mechanism for maintaining an initialized
multi-threaded component object model apartment.

When working with mutli-threaded apartments, COM requires at least one
thread to be initialized, otherwise COM-allocated resources may be released
prematurely. This poses a challenge in Go, which can have many goroutines
running in parallel with weak thread affinity.

The comshim package provides a solution to this problem by maintaining
a single thread-locked goroutine that has been initialized for
multi-threaded COM use via a call to CoIntializeEx. A reference counter is
used to determine the ongoing need for the shim to stay in place. Once the
counter reaches 0, the thread is released and COM may be deinitialized.

The comshim package is designed to allow COM-based libraries to hide the
threading requirements of COM from the user. COM interfaces can be hidden
behind idomatic Go structures that increment the counter with calls to
NewType() and decrement the counter with calls to Type.Close(). To see
how this is done, take a look at the WrapperUsage example.

Example Usage
====

```
package main

import "github.com/scjalliance/comshim"

func main() {
	// This ensures that at least one thread maintains an initialized
	// multi-threaded COM apartment.
	comshim.Add(1)

	// After we're done using COM the thread will be released.
	defer comshim.Done()

	// Do COM things here
}
```