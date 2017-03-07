package msg

type message struct {
	alivetime string
	body      string
}

func newMessage(alivetime, body string) *message {
	return &message{
		alivetime: alivetime,
		body:      body,
	}
}
