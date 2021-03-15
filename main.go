package main

import (
	"ethersim/common"
	"ethersim/media"
	"ethersim/tap"
	"fmt"
	"os"
	"strings"
)

func main() {
	conf, err := common.ArgParse(os.Args[1:])
	if err != nil {
		fmt.Println("Usage: ethersim [OPTIONS]")
		fmt.Println("\t-s: Start ethersim server / emulated media")
		fmt.Println("\t-c [ID]: Connect to session [ID]")
		return
	}

	if conf.Serving {
		sessionStart := strings.LastIndex(conf.Name, "/") + 1
		sessionID := conf.Name[sessionStart:]
		fmt.Printf("Good to go! Your ether ID is %s.\n", sessionID)
		media.ServeCon(&conf)

	} else {
		tap.JoinSession(&conf)
	}
}
