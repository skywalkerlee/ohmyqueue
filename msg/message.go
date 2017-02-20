package msg

type Header struct {
	Len      int
	Deadline int64
}

type Msg struct {
	Header Header
	Body   string
}
