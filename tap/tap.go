package tap

import (
	"ethersim/common"
	"fmt"
	"log"
	"net"
	"time"
	"unicode/utf8"
)

// JoinSession sets up a connection with
// an ethersim server, performing multithreaded
// IO with the UDS on the other end.
func JoinSession(c *common.Conf) {
	con, err := net.Dial("unix", c.Name)
	if err != nil {
		log.Fatalf("Fatal error establishing connection: %s", err.Error())
	}

	go doOutput(&con)
	doInput(&con)
}

func doOutput(c *net.Conn) {
	b := make([]byte, 4)
	isReceiving := false

	for {
		_, err := (*c).Read(b)
		if err != nil {
			(*c).Close()
			return
		}

		r, _ := utf8.DecodeRune(b)
		if r == utf8.RuneError {
			continue
		}

		if !isReceiving {
			fmt.Print("Receiving: ")
			isReceiving = true
		}

		if r == '✓' {
			isReceiving = false
		} else {
			fmt.Printf("%c", r)
		}

	}
}

func doInput(c *net.Conn) {
	for {
		str := readMsgChars()
		if str == "\n" {
			continue
		}

		fmt.Print("Sending... ")
		b := make([]byte, 4)

		for _, r := range str {
			fmt.Printf("%c", r)
			utf8.EncodeRune(b, r)

			_, err := (*c).Write(b)
			if err != nil {
				(*c).Close()
				return
			}
			time.Sleep(100 * time.Millisecond)
		}

		utf8.EncodeRune(b, '✓')
		_, err := (*c).Write(b)
		if err != nil {
			(*c).Close()
			return
		}
	}
}
