package msg

type Topics map[string]*topic

func NewTopics() Topics {
	topics := make(map[string]*topic)
	return topics
}

func (topics Topics) AddTopic(name string) {
	topics[name] = newTopic(name)
	topics[name].load()
	go topics[name].clean()
}

func (topics Topics) Put(topic, alivetime, body string, offset ...string) {
	topics[topic].put(alivetime, body, offset...)
}

func (topics Topics) Get(topic, offset string) *message {
	return topics[topic].get(offset)
}
