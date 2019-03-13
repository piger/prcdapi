package main

import (
	"flag"
	"log"
	"os"

	"github.com/piger/prcdapi"
)

var (
	address = flag.String("address", "127.0.0.1:30666", "Specify the bind address (default: 127.0.0.1:30666)")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	dbDir := args[0]
	grimoire, err := prcdapi.LoadPrcdDir(dbDir)
	if err != nil {
		log.Fatal(err)
	}

	s := prcdapi.NewServer(grimoire)
	s.Serve(*address)
}
