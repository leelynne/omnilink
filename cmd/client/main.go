package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/leelynne/omnilink/omni/home"
)

func main() {
	logger := log.New(os.Stdout, "omni: ", log.LstdFlags)
	var endpoint string
	var key string
	flag.StringVar(&endpoint, "endpoint", "", "endpoint to connect to.")
	flag.StringVar(&key, "key", "", "client key")
	flag.Parse()

	h, err := home.New(logger, endpoint, key)
	if err != nil {
		panic(err)
	}
	fmt.Println(h)
}
