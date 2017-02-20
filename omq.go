package main

import (
	"github.com/ohmq/ohmyqueue/msg"
	"github.com/ohmq/ohmyqueue/server"
)

func main() {
	msgs := msg.NewMsgs()
	broker := server.NewBroker(1, msgs)
	broker.Start()
}
