package main

import (
	"flag"
	"fmt"

	"github.com/leelynne/omnilink/omni"
)

func main() {
	var endpoint string
	var key string
	flag.StringVar(&endpoint, "endpoint", "", "endpoint to connect to.")
	flag.StringVar(&key, "key", "", "client key")
	flag.Parse()
	c, err := omni.NewClient(fmt.Sprintf("%s:4369", endpoint), key)

	if err != nil {
		panic(err)
	}

	si, err := c.GetSystemInformation()
	if err != nil {
		panic(err)
	}
	fmt.Printf("SysINFO %+v\n", si)
	fmt.Printf("SysINFO %s\n", string(si.LocalPhoneNumber[:]))
}
