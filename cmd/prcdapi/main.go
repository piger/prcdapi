package main

import (
	"flag"
	"log"

	"github.com/piger/prcdapi"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Missing db directory")
	}
	dbDir := args[0]

	grimoire, err := prcdapi.LoadPrcdDir(dbDir)
	if err != nil {
		log.Fatal(err)
	}

	s := prcdapi.NewServer(grimoire)
	s.Serve()
}
