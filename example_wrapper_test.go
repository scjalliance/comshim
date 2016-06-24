package comshim_test

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

func Example_wrapperUsage() {
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
