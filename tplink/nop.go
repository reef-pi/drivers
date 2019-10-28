package tplink

import (
	"io"
	"time"
)

type nopConn struct {
	read   bool
	Buffer []byte
}

func (c *nopConn) Close() error { return nil }
func (c *nopConn) Read(buf []byte) (int, error) {
	if c.read {
		return 0, io.EOF
	}
	copy(buf, c.Buffer)
	c.read = true
	return len(buf), nil
}
func (c *nopConn) SetDeadline(_ time.Time) error { return nil }
func (c *nopConn) Write(_ []byte) (int, error)   { return 0, nil }

type nop struct {
	read bool
	conn *nopConn
}

func (n *nop) Buffer(b []byte) {
	n.conn = &nopConn{
		Buffer: b,
	}
}

func (n *nop) Factory(_, _ string, _ time.Duration) (Conn, error) {
	return n.conn, nil
}

func NewNop() *nop {
	return &nop{
		conn: &nopConn{
			Buffer: []byte(`{}`),
		},
	}
}
