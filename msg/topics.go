package msg

type Topics map[string]*topic

func NewTopics() Topics {
	topics := make(map[string]*topic)
	return topics
}

func (topics Topics) AddTopic(name string) {
	topics[name] = newTopic(name)
}

func (topics Topics) Put(topic, alivetime, body string) {
	topics[topic].put(alivetime, body)
}

func (topics Topics) Get(topic, offset string) *message {
	return topics[topic].get(offset)
}
