package tap

import (
	"bufio"
	"ethersim/common"
	"fmt"
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

func doOutput(c *net.Conn, s chan<- bool, wc chan int) {
	b := make([]byte, 4)
	isReceiving := false

	for {
		if delay := tryExpBackoff(wc); delay > 0 {
			// Eat all the datas, to be safe
			tmpbuf := make([]byte, delay*10+1)
			(*c).Read(tmpbuf)
			isReceiving = false
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
			fmt.Print("Receiving: ")
			isReceiving = true
			s <- true
		}

		if r == '✓' {
			isReceiving = false
			s <- false
		} else if r == '☠' {
			isReceiving = false
			s <- true
			s <- false
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

			fmt.Print("Sending... ")
			b := make([]byte, 4)
			s <- true

			for _, r := range str {
				if i := tryExpBackoff(wc); i > 0 {
					break
				}
				fmt.Printf("%c", r)
				utf8.EncodeRune(b, r)

				_, err := (*c).Write(b)
				if err != nil {
					(*c).Close()
					return
				}
				time.Sleep(time.Duration(common.MsgInterval) * time.Millisecond)
			}

			utf8.EncodeRune(b, '✓')
			_, err := (*c).Write(b)
			s <- false
			if err != nil {
				(*c).Close()
				return
			}
		}
	}
}
