package server

import (
	"net"

	"github.com/ik5/distributed_test/log"
	"github.com/ik5/distributed_test/types"
)

// Handler is a callback for any given server
type Handler = func(conn net.Conn, buf []byte)

// Listen listen to an address and port using a UDP, and executing a callback
// when a callback when content arrives
func Listen(addr string, handler Handler) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}
		buf := make([]byte, types.MaxBufferSize)
		n, err := conn.Read(buf)
		if err != nil {
			log.Logger.Tracef("Unable to read buf: %s", err)
			continue
		}
		go handler(conn, buf[:n])
	}
}
