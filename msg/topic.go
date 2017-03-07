package msg

import (
	"strconv"
)

type topic map[string]*message

func newTopic() topic {
	topic := make(map[string]*message)
	return topic
}

func (topic topic) put(alivetime, body string) {
	topic[strconv.Itoa(len(topic))] = newMessage(alivetime, body)
}

func (topic topic) get(offset string) *message {
	return topic[offset]
}
