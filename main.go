package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/leelynne/omnilink/omni"
)

func main() {
	logger := log.New(os.Stdout, "omni: ", log.LstdFlags)
	var endpoint string
	var key string
	flag.StringVar(&endpoint, "endpoint", "", "endpoint to connect to.")
	flag.StringVar(&key, "key", "", "client key")
	flag.Parse()
	c, err := omni.NewClient(fmt.Sprintf("%s:4369", endpoint), key)

	if err != nil {
		panic(err)
	}

	logger.Printf("Connected!")
	si, err := c.GetSystemInformation()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SysINFO %+v\n", si)
	fmt.Printf("Phone '%s'\n", string(si.LocalPhoneNumber[:]))

	st, err := c.GetSystemStatus()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SysSTATUS %+v\n", st)

}
