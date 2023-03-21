package reply

var pongBytes = []byte("+PONG\r\n")
var okBytes = []byte("+OK\r\n")
var nullBulkBytes = []byte("$-1\r\n")
var emptyMultiBulkBytes = []byte("*0\r\n")
var noBytes = []byte("")

var thePongReply = new(PongReply)
var theOkReply = new(OkReply)
var theNullBulkReply = new(nullBulkReply)
var theEmptyMultiBulkReply = new(emptyMultiBulkReply)
var theNoReply = new(NoReply)

/******************************************************************************/

type PongReply struct {
}

func NewPongReply() *PongReply {
	return &PongReply{}
}

func (p *PongReply) ToBytes() []byte {
	return pongBytes
}

/******************************************************************************/

type OkReply struct {
}

func (t *OkReply) ToBytes() []byte {
	return okBytes
}

func NewOkReply() *OkReply {
	return &OkReply{}
}

/******************************************************************************/

type nullBulkReply struct {
}

func NewNullBulkReply() *nullBulkReply {
	return &nullBulkReply{}
}

func (n *nullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

/******************************************************************************/

type emptyMultiBulkReply struct {
}

func NewEmptyMultiBulkReply() *emptyMultiBulkReply {
	return &emptyMultiBulkReply{}
}

func (e *emptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

/******************************************************************************/

type NoReply struct {
}

func (n NoReply) ToBytes() []byte {
	return noBytes
}

func NewNoReply() *NoReply {
	return &NoReply{}
}
