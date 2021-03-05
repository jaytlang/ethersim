package media

import (
	"log"
	"net"
	"unicode/utf8"
)

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
func input(con *net.Conn, c chan<- rune) {
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
		c <- r
	}
}

func broadcast(dc <-chan rune, cc <-chan *net.Conn) {
	cons := []*net.Conn{}

	for {
		select {

		// If there's a new connection, add
		// it to the list.
		case newcon := <-cc:
			cons = append(cons, newcon)

		// If there's new data, send it along
		case r := <-dc:
			bytes := make([]byte, 4)
			utf8.EncodeRune(bytes, r)

			for i, c := range cons {
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
func ServeCon(nm string) {
	ln, err := net.Listen("unix", nm)
	if err != nil {
		log.Fatalf("Fatal error starting ethersim server: %s", err.Error())
	}

	dc := make(chan rune)
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
	}
}
