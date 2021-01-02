package app

import (
	"io"
	"net"
	"os"

	"github.com/kdevb0x/or-maybe-more/or-maybe-more/node"
)

var Log = newLogger(os.Stderr)

type connListener struct {
	addr string
	conn *net.IPConn
}

func (cl *connListener) Accept() (net.Conn, error) {
	return cl.conn, nil
}

func (cl *connListener) Addr() net.Addr {
	return cl.conn.LocalAddr()
}

func (cl *connListener) Close() error {
	return cl.conn.Close()
}

type LocationInfo = node.LocInfo

// Transport is an abstract net.Conn supporting our custom protocols.
type Transport interface {
	net.Conn
}

type netchanTransport struct {
	local, remote io.ReadWriteCloser
}

func (t *transport) Write(p []byte) (int, error) {
	return t.remote.Write(p)
}
