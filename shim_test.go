package comshim_test

import (
	"runtime"
	"sync"
	"testing"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"

	"github.com/NozomiNetworks/go-comshim"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentShims(t *testing.T) {
	comLoader := comshim.NewLoader()
	defer comLoader.Wait()
	var maxRounds int
	if testing.Short() {
		maxRounds = 64
	} else {
		maxRounds = 256
	}

	// Vary the number of threads
	for procs := 1; procs < 11; procs++ {
		runtime.GOMAXPROCS(procs)

		// Vary the number of shims
		for rounds := 1; rounds <= maxRounds; rounds *= 2 {
			wg := sync.WaitGroup{}
			for i := 0; i < rounds; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()

					loadErr := comLoader.Load()
					assert.Nil(t, loadErr)
					defer comLoader.Unload()

					obj, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
					if err != nil {
						t.Error(err)
					} else {
						defer obj.Release()
					}
				}(i)
			}
			wg.Wait()
		}
	}
}

func TestConcurrentCoInitializeDoesNotPanic(t *testing.T) {
	comLoader := comshim.NewLoader()
	defer comLoader.Wait()
	var maxRounds int
	if testing.Short() {
		maxRounds = 64
	} else {
		maxRounds = 256
	}

	// Vary the number of threads
	for procs := 1; procs < 11; procs++ {
		runtime.GOMAXPROCS(procs)

		// Vary the number of shims
		for rounds := 1; rounds <= maxRounds; rounds *= 2 {
			wg := sync.WaitGroup{}
			for i := 0; i < rounds; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_ = comLoader.Load()
					defer comLoader.Unload()
				}()

				wg.Add(1)
				go func() {
					defer wg.Done()
					_ = ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
				}()
			}
			wg.Wait()
		}
	}
}
