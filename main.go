package main

import (
	"log"

	core "github.com/andryanduta/taxi-fare/core"
)

func main() {
	// init config
	_, err := core.InitMainConfig()
	if err != nil {
		log.Fatalln(err)
	}
	// init log

	// init domain

	// handler start
}
