package tap

import (
	"time"
)

// Currently disabled, returns 2 always
func mkInterval(ctr int) int {
	return 2
}

// Given the con and the sleep channel, sleep
// if we need to. Returns true if we actually did.
func tryExpBackoff(wc <-chan int) int {
	select {
	case delay := <-wc:
		// Sleep for the correct interval...
		time.Sleep(time.Duration(delay) * time.Second)
		return delay
	default:
		return 0
	}
}
