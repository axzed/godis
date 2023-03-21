package resp

// Reply is the interface that wraps the basic ToBytes method.
type Reply interface {
	ToBytes() []byte
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}
