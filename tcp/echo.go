package tcp

import (
	"bufio"
	"context"
	"godis/lib/logger"
	atomic2 "godis/lib/sync/atomic"
	"godis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

// EchoClient is a client for echo server
// Conn is the connection (socket)
// Waiting is used to wait for the client to close
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close closes the client
func (e *EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second) // wait for 10 seconds to close the client
	_ = e.Conn.Close()                          // close the connection
	return nil
}

// EchoHandler is a handler for echo server
// activeConn is a map to store all active connections, using sync.Map to make it thread-safe
// closing is used to indicate whether the server is closing
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic2.Boolean
}

func NewHandler() *EchoHandler {
	return &EchoHandler{}
}

// Handle handles the connection
func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	// close the connection when the server is closing
	if handler.closing.Get() {
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	handler.activeConn.Store(client, struct{}{}) // store the connection in the map
	reader := bufio.NewReader(conn)              // use buffer to read data from the client
	for {
		msg, err := reader.ReadString('\n') // \n is the delimiter of the message
		if err != nil {
			if err == io.EOF {
				logger.Info("Connection close")
				handler.activeConn.Delete(client) // delete the connection from the map
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1) // add 1 to Waiting when you start a goroutine to write data to the client
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done() // minus 1 from Waiting when you finish writing data to the client
	}
}

// Close closes the handler
func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	handler.closing.Set(true) // set the closing flag to true means the server is closing
	handler.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Conn.Close()
		return true // return true to continue the iteration
	})
	return nil
}
