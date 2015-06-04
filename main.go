package main

import "fmt"

func main() {
	c := Client{
		Addr: "199.231.243.13:4369",
	}
	err := c.Connect()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected!")
	m := msg{
		SeqNum:  1,
		MsgType: CLIENT_REQ_NEW_SESSION,
	}
	err = c.Send(m)
	if err != nil {
		panic(err)
	}
	resp, err := c.Receive()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Msg %+v", resp)
}
