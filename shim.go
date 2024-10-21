package comshim

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/go-ole/go-ole"
)

var (
	ErrNegativeCounter    = errors.New("COM already unloaded")
	ErrAlreadyInitialized = errors.New("COM thread has already been initialized")
)

// Loader maintains CoInitializeEx as long as required
// https://learn.microsoft.com/en-us/windows/win32/api/combaseapi/nf-combaseapi-coinitializeex
type Loader struct {
	startAccess  sync.Mutex
	loaded       bool
	signalAccess sync.Mutex
	signal       sync.Cond
	workTotal    atomic.Int64 // https://pkg.go.dev/sync/atomic#pkg-note-BUG
	wg           sync.WaitGroup
}

func NewLoader() *Loader {
	shim := Loader{}
	shim.signal.L = &shim.signalAccess
	return &shim
}

// Load loads COM, call Unload when COM is no longer required.
// If Load returns ErrAlreadyInitialized subsequent COM calls might still succeed,
// as long as the other location where COM was loaded maintains it.
func (s *Loader) Load() error {
	s.startAccess.Lock()
	defer s.startAccess.Unlock()
	s.requiredBy(1)
	if s.loaded {
		return nil // already loaded
	}

	err := s.start()
	if err != nil {
		return err
	}

	s.loaded = true
	return nil
}

func (s *Loader) Unload() {
	s.requiredBy(-1)
}

func (s *Loader) requiredBy(delta int64) {
	s.signalAccess.Lock()
	defer s.signalAccess.Unlock()
	value := s.workTotal.Add(delta)
	s.signal.Signal()
	if value < 0 { // invalid usage
		panic(ErrNegativeCounter)
	}
}

func (s *Loader) start() error {
	init := make(chan error)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
		if err != nil {
			switch err.(*ole.OleError).Code() {
			case 0x00000001: // windows.S_FALSE
				ole.CoUninitialize()
				init <- ErrAlreadyInitialized
			default:
				init <- err
			}
			close(init)
			return
		}

		close(init)

		{ // work until no longer required
			s.signalAccess.Lock()
			for s.workTotal.Load() > 0 {
				s.signal.Wait()
			}
			s.loaded = false
			ole.CoUninitialize()
			s.signalAccess.Unlock()
		}
	}()

	return <-init
}

func (s *Loader) Wait() {
	s.startAccess.Lock()
	defer s.startAccess.Unlock()
	s.wg.Wait()
}
