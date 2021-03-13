package media

import (
	"ethersim/common"
	"log"
	"net"
	"os"
	"unicode/utf8"
)

type traceData struct {
	r   rune
	con *net.Conn
}

// The following are input and broadcast
// operations for our one-to-many setup here.
// Basically, if one client we're connected to
// says something, we talk to ALL of them.
//
// Notably, this doesn't fit the 033 personality well
// because it uses CSP concurrency rather than
// the l.l. locking things we teach...maybe looking
// at this code isn't too instructive. However, for this
// use case, CSP is REALLY REALLY nice.
func input(con *net.Conn, c chan<- traceData) {
	for {
		// Receive a new message, probably
		// in the form of a 4-byte rune
		msg := make([]byte, 4)
		_, err := (*con).Read(msg)
		if err != nil {
			(*con).Close()
			return
		}

		// Decode said message, which has
		// the added benefit of error checking
		r, _ := utf8.DecodeRune(msg)
		if r == utf8.RuneError {
			continue
		}

		// Send it along!
		c <- traceData{r: r, con: con}
	}
}

func broadcast(dc <-chan traceData, cc <-chan *net.Conn) {
	cons := []*net.Conn{}

	for {
		select {

		// If there's a new connection, add
		// it to the list.
		case newcon := <-cc:
			cons = append(cons, newcon)

		// If there's new data, send it along
		case td := <-dc:
			bytes := make([]byte, 4)
			utf8.EncodeRune(bytes, td.r)

			for i, c := range cons {
				if c == td.con {
					continue
				}
				_, err := (*c).Write(bytes)
				if err != nil {
					(*c).Close()
					cons = append(cons[:i], cons[i+1:]...)
				}
			}
		}
	}
}

// ServeCon randomly generates a new connection
// identifier, and begins serving it. This continues
// forever or so -- until the caller tells us to stop.
func ServeCon(c *common.Conf) {
	ln, err := net.Listen("unix", (*c).Name)
	if err != nil {
		log.Fatalf("Fatal error starting ethersim server: %s", err.Error())
	}

	if err := os.Chmod((*c).Name, 0777); err != nil {
		log.Fatalf("Error setting permissions: %s\n", err.Error())
	}

	dc := make(chan traceData)
	cc := make(chan *net.Conn)

	// Start the broadcaster
	go broadcast(dc, cc)

	for {
		// Accept a new connection
		con, err := ln.Accept()
		if err != nil {
			log.Fatalf("Error accepting new connection: %s", err.Error())
		}

		// Aggregate the new connection and spin
		// an input for it
		cc <- &con
		go input(&con, dc)
	}
}
