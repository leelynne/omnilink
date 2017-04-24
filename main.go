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
		fmt.Printf("%+v\n", err)
		panic(err)
	}

	logger.Printf("Connected!")
	si, err := c.GetSystemInformation()
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	fmt.Printf("SysINFO %+v\n", si)
	fmt.Printf("Phone '%s'\n", string(si.LocalPhoneNumber[:]))

	st, err := c.GetSystemStatus()
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	fmt.Printf("SysSTATUS %+v\n", st)

	tr, err := c.GetSystemTroubles()
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	fmt.Printf("SysTroubles %+v\n", tr)

	ftr, err := c.GetSystemFeatures()
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	fmt.Printf("SysFeatures %+v\n", ftr)

	form, err := c.GetSystemFormats()
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	fmt.Printf("SysFormats %+v\n", form)

	otc, err := c.GetObjectTypeCapacity(omni.Area)
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	fmt.Printf("Capacity %+v\n", otc)

	oprop, err := c.GetObjectProperties()
	if err != nil {
		fmt.Printf("%+v\n", err)
		panic(err)
	}
	fmt.Printf("Property %+v\n", oprop)

}
