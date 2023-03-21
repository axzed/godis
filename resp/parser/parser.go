package parser

import (
	"bufio"
	"errors"
	"godis/interface/resp"
	"io"
)

// Payload use the reply
type Payload struct {
	Data resp.Reply
	Err  error
}

// readState is the state of the parser
type readState struct {
	readingMultiLine  bool
	expectedArgsCount int
	msgType           byte
	args              [][]byte
	bulkLen           int64
}

// finished means the parser has finished parsing a command
func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount
}

// ParseStream parses the stream byte by byte
// use goroutine to parse the stream in parse0
// TCP conn use this method to parse the stream
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

// parse0 is the core method of the parser
// it parses the stream byte by byte
// and sends the parsed result to the channel
// ParseStream will async call this method
func parse0(reader io.Reader, ch chan<- *Payload) {

}

// readLine reads a line from the reader
// []byte is the line, bool is the flag whether the parser has IO error
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	// 1. \r\n slice
	// 2. if read the $ before, count the char strictly
	var msg []byte
	var err error
	if state.bulkLen == 0 { // \r\n slice
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
	} else { // 2. if read the $ before, count the char strictly
		msg = make([]byte, state.bulkLen+2)
		_, err := io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
		// next line
		state.bulkLen = 0
	}
	return msg, false, nil
}