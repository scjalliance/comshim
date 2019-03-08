package comshim_test

import "github.com/scjalliance/comshim"

func Example_globalUsage() {
	// This ensures that at least one thread maintains an initialized
	// multi-threaded COM apartment.
	comshim.Add(1)

	// After we're done using COM the thread will be released.
	defer comshim.Done()

	// Do COM things here
}
