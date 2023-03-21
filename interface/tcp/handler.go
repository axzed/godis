package tcp

import (
	"context"
	"net"
)

// Handler TCP 业务逻辑处理接口
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
