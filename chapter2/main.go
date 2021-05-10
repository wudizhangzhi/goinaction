package main

import (
	_ "chapter2/matchers"
	"chapter2/search"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	search.Run(os.Args[1])
}
