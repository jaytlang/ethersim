package tap

import (
	"ethersim/common"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

// JoinSession sets up a connection with
// an ethersim server, performing multithreaded
// IO with the UDS on the other end.
func JoinSession(c *common.Conf) {
	rand.Seed(time.Now().UnixNano())

	con, err := net.Dial("unix", c.Name)
	if err != nil {
		log.Fatalf("Fatal error establishing connection: %s", err.Error())
	}

	inC := make(chan bool)
	outC := make(chan bool)
	inS := make(chan int, 1)
	outS := make(chan int, 1)

	go doOutput(&con, outC, outS)
	go doInput(&con, inC, inS)

	first, last := true, true
	ctr := 0

	for n := range merge(inC, outC) {
		if first {
			first = false
			last = n
			continue
		}

		if n == last && n == true {
			ctr++
			interval := mkInterval(ctr)

			outS <- interval
			inS <- interval

			fmt.Println("\n****** A conflict occurred! ******")
			fmt.Println("Someone tried to send and receive at the same time")
			fmt.Println("Applying backoff with interval", interval)
			fmt.Println("***********************************")
		} else {
			ctr = 0
		}
		last = n
	}

}
