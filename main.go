package main

import (
	"github.com/crazy-me/ops_dial/logic"
	"log"
	"os"
)

var commandArgs []string

func main() {
	commandArgs = os.Args
	if 2 != len(commandArgs) {
		log.Fatalln("Inaccurate operating parameters")
	}
	logic.Run(commandArgs[1])
}
