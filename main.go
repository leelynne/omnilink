package main

import "github.com/leelynne/omnilink/omni"

func main() {
	_, err := omni.NewClient("199.231.243.13:4369", "")

	if err != nil {
		panic(err)
	}
}
