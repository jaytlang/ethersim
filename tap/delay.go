package tap

import (
	"math"
	"math/rand"
	"time"
)

// Slightly modified for effect
func mkInterval(ctr int) int {
	rval := rand.Intn(ctr + 1)
	return int(math.Pow(2, float64(rval)) - 1)
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
