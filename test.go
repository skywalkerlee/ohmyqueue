package main

import (
	"log"

	"github.com/ohmq/ohmyqueue/msg"
)

func main() {
	msgs := msg.NewMsgs()
	msgs.Update([]string{"test1", "test2"})
	msg := msg.Msg{Header: msg.Header{Len: 10, Deadline: 10}, Body: "aaa"}
	log.Printf("%#v", msg)
	msgs.Put("test1", msg)
}
