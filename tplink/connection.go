package tplink

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"time"
)

const (
	_timeOut    = 2 * time.Second
	_buffLength = 512
)

type Conn interface {
	Close() error
	SetDeadline(time.Time) error
	Write([]byte) (int, error)
	Read([]byte) (int, error)
}

type ConnectionFactory func(string, string, time.Duration) (Conn, error)

var TCPConnFactory = func(proto, addr string, t time.Duration) (Conn, error) {
	return net.DialTimeout(proto, addr, t)
}

type cmd struct {
	cf   ConnectionFactory
	addr string
}

func (c *cmd) Execute(command interface{}, pResult bool) ([]byte, error) {
	payload, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}
	conn, err := c.cf("tcp", c.addr, _timeOut)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(_timeOut)); err != nil {
		return nil, err
	}
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(payload)))
	bs := append(header, autokeyEncrypt(payload)...)
	_, err = conn.Write(bs)
	if err != nil {
		return nil, err
	}
	if !pResult {
		return []byte{}, nil
	}
	if _, err := conn.Read(header); err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, rErr := io.Copy(buf, conn)
	resp := buf.Bytes()
	if len(resp) == 0 && rErr != nil {
		return nil, rErr
	}
	return autokeyDecrypt(resp), nil
}
