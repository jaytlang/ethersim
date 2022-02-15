package tap

import (
	"bufio"
	"ethersim/common"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
	"unicode/utf8"
)

func readMsgChars() string {
	fmt.Print("Type your message here: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

func doOutput(c *net.Conn, s chan<- bool, wc <-chan int) {
	b := make([]byte, 4)
	isReceiving := false

	for {
		if delay := tryExpBackoff(wc); delay > 0 {
			// Eat all the datas, to be safe
			// Delay is in terms of seconds, and we send
			// one message every 100 milliseconds, so this
			// works out
			tmpbuf := make([]byte, delay*10)
			(*c).Read(tmpbuf)
			continue
		}

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
			fmt.Print("\nReceiving: ")
			isReceiving = true
			s <- true
		}

		if r == '✓' {
			fmt.Print("Type your message here: ")
			isReceiving = false
			s <- false
		} else if r == '✖' {
			isReceiving = false
			s <- false
			fmt.Printf("\n***Note: CRC32 checksum failed on this message!***")
		} else {
			fmt.Printf("%c", r)
		}

	}
}

func doInput(c *net.Conn, s chan<- bool, wc <-chan int) {
	for {
		// If we are here, trash the expbackoff
		select {
		case <-wc:
			continue
		default:

			str := readMsgChars()
			if str == "\n" {
				continue
			}

			willFail := false

			fmt.Print("Sending... ")
			b := make([]byte, 4)
			s <- true

			for _, r := range str {
				if i := tryExpBackoff(wc); i > 0 {
					break
				}

				if r == '`' {
					willFail = true
					continue
				}

				fmt.Printf("%c", r)
				if willFail {
					r += rune(rand.Intn(5)) - rune(rand.Intn(10))
				}
				utf8.EncodeRune(b, r)

				_, err := (*c).Write(b)
				if err != nil {
					(*c).Close()
					return
				}
				time.Sleep(time.Duration(common.MsgInterval) * time.Millisecond)
			}

			if willFail {
				utf8.EncodeRune(b, '✖')
			} else {
				utf8.EncodeRune(b, '✓')
			}
			_, err := (*c).Write(b)
			s <- false
			if err != nil {
				(*c).Close()
				return
			}
		}
	}
}
