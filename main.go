package main

import (
	"github.com/romanornr/cyberchain/client"
	"github.com/romanornr/cyberchain/server"
)

func main() {
	server.Start()
	cl := client.GetInstance()
	cl.Connect(2)
}
