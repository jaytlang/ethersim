package main

import (
	"ethersim/conf"
	"ethersim/media"
	"log"
)

func main() {
	conf, err := conf.MkConfig()
	if err != nil {
		log.Fatal("Invalid arguments!")
	}

	nm := conf.MkName()
	media.ServeCon(nm)
}
