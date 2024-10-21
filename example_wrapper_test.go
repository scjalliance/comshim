package comshim_test

import (
	"sync"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// Object wraps a COM interface in a way that is safe for multi-threaded access.
// In this example it wraps IUnknown.
type Object struct {
	m         sync.Mutex
	iface     *ole.IUnknown
	comLoader *comshim.Loader
}

// NewObject creates a new object. Be sure to document the need to call Close().
func NewObject(comLoader *comshim.Loader) (*Object, error) {
	_ = comLoader.Load()
	iunknown, err := oleutil.CreateObject("Excel.Application")
	if err != nil {
		comLoader.Unload()
		return nil, err
	}
	return &Object{iface: iunknown, comLoader: comLoader}, nil
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
	o.comLoader.Unload()
}

// Foo performs some action using the object's COM interface.
func (o *Object) Foo() {
	o.m.Lock()
	defer o.m.Unlock()

	// Make use of o.iface
}

func Example_wrapperUsage() {
	comLoader := comshim.NewLoader()
	defer comLoader.Wait()
	obj1, err := NewObject(comLoader) // Create an object
	if err != nil {
		panic(err)
	}
	defer obj1.Close() // Be sure to close the object when finished

	obj2, err := NewObject(comLoader) // Create a second object
	if err != nil {
		panic(err)
	}
	defer obj2.Close() // Be sure to close it too

	// Work with obj1 and obj2
}
