package msg

import (
	"strconv"
	"sync"
)

type topic struct {
	mutex *sync.Mutex
	msg   map[string]*message
}

func newTopic() *topic {
	msg := make(map[string]*message)
	return &topic{
		mutex: new(sync.Mutex),
		msg:   msg,
	}
}

func (topic *topic) put(alivetime, body string) {
	topic.mutex.Lock()
	topic.msg[strconv.Itoa(len(topic.msg))] = newMessage(alivetime, body)
	topic.mutex.Unlock()
}

func (topic *topic) get(offset string) *message {
	return topic.msg[offset]
}
