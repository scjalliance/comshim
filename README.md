comshim [![GoDoc](https://godoc.org/github.com/scjalliance/comshim?status.svg)](https://godoc.org/github.com/scjalliance/comshim)
====

The comshim package provides a mechanism for maintaining an initialized
multi-threaded component object model apartment.

When working with mutli-threaded apartments, COM requires at least one
thread to be initialized, otherwise COM-allocated resources may be released
prematurely. This poses a challenge in Go, which can have many goroutines
running in parallel with weak thread affinity.

The `comshim` package provides a solution to this problem by maintaining
a single thread-locked goroutine that has been initialized for
multi-threaded COM use via a call to `CoIntializeEx`. A reference counter is
used to determine the ongoing need for the shim to stay in place. Once the
counter reaches 0, the thread is released and COM may be deinitialized.

The comshim package is designed to allow COM-based libraries to hide the
threading requirements of COM from the user. COM interfaces can be hidden
behind idomatic Go structures that increment the counter with calls to
`NewType()` and decrement the counter with calls to `Type.Close()`. To see
how this is done, take a look at the WrapperUsage example.

Global Example
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

Wrapper Example
====

```
package main

import (
	"sync"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/scjalliance/comshim"
)

// Object wraps a COM interface in a way that is safe for multi-threaded access.
// In this example it wraps IUnknown.
type Object struct {
	m     sync.Mutex
	iface *ole.IUnknown
}

// NewObject creates a new object. Be sure to document the need to call Close().
func NewObject() (*Object, error) {
	comshim.Add(1)
	iunknown, err := oleutil.CreateObject("Excel.Application")
	if err != nil {
		comshim.Done()
		return nil, err
	}
	return &Object{iface: iunknown}, nil
}

// Close releases any resources used by the object.
func (o *Object) Close() {
	o.m.Lock()
	defer o.m.Unlock()
	if o.iface == nil {
		return // Already closed
	}
	o.iface.Release()
	o.iface = nil
	comshim.Done()
}

// Foo performs some action using the object's COM interface.
func (o *Object) Foo() {
	o.m.Lock()
	defer o.m.Unlock()

	// Make use of o.iface
}

func main() {
	obj1, err := NewObject() // Create an object
	if err != nil {
		panic(err)
	}
	defer obj1.Close() // Be sure to close the object when finished

	obj2, err := NewObject() // Create a second object
	if err != nil {
		panic(err)
	}
	defer obj2.Close() // Be sure to close it too

	// Work with obj1 and obj2
}
```