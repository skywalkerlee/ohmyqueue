package msg

type Header struct {
	len      int
	deadline int64
}

type Msg struct {
	header Header
	body   string
}
