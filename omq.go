package main

import (
	"github.com/skywalkerlee/ohmyqueue/server"
)

func main() {
	broker := server.NewBroker(1)
	broker.Start()
}
