package tap

import "sync"

// Function to merge a handful of channels.
// Taken from the golang blog -- all credit goes
// there!! We'll comment here to be sure
// that we understand it, once we verify it works.
func merge(cs ...<-chan bool) <-chan bool {
	var wg sync.WaitGroup
	out := make(chan bool)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan bool) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
