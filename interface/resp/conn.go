package resp

// Connection is the interface that wraps the basic Write method.
type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
