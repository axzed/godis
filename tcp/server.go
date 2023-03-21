package tcp

import (
	"context"
	"godis/interface/tcp"
	"godis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Config TCP server config
type Config struct {
	Address string
}

// ListenAndServeWithSignal starts a TCP server with signal
// create a goroutine to listen on the signal channel
// when signal received, close the closeChan
// when closeChan closed, stop the server
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	// Notify the signal channel when the process receives SIGHUP, SIGQUIT, SIGTERM, SIGINT
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		// when the process receives SIGHUP, SIGQUIT, SIGTERM, SIGINT, close the closeChan
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			// send a signal to closeChan to stop the server
			closeChan <- struct{}{}
		}
	}()
	listener, err := net.Listen("tcp", cfg.Address) // create the listener (TCP socket)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	ListenAndServe(listener, handler, closeChan) // start the server
	return nil
}

// ListenAndServe starts a TCP server
// handle connections in goroutines
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	go func() {
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()
	}()

	// close the listener and handler when the server stops
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()
	var waitDone sync.WaitGroup

	// continue to accept connections until closeChan closed
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted link")
		// add 1 to waitDone when you start a goroutine to handle the connection
		waitDone.Add(1)
		// goroutine-per-connection
		go func() {
			// when the connection is handled, minus 1 from waitDone
			defer func() {
				waitDone.Done()
			}()
			// handle the connection
			handler.Handle(ctx, conn)
		}()
	}

	// when you break the for loop it means the TCP server failed to accept a connections
	// but the server still needs to wait for another connections to be handled
	waitDone.Wait()
}
