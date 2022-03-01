package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	commandArgs []string
	fileName    string
)

func main() {
	initArgs()
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(fileName, "Open failed:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		accountSlice := strings.Split(scanner.Text(), "|")
		fmt.Println(accountSlice)
	}

}

func initArgs() {
	commandArgs = os.Args
	if 2 != len(commandArgs) {
		log.Fatalln("Inaccurate operating parameters")
	}
	fileName = commandArgs[1]
}
