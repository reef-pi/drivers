package tplink

import (
	"encoding/binary"
	"net"
	"sync"
	"time"
)

var (
	_timeOut = 5 * time.Second
	_cmdOn   = []byte(`{"system":{"set_relay_state":{"state":1}}}`)
	_cmdOff  = []byte(`{"system":{"set_relay_state":{"state":0}}}`)
	_cmdInfo = []byte(`{"system":{"get_sysinfo":{}}}`)
)

type Conn interface {
	Close() error
	SetDeadline(time.Time) error
	Write([]byte) (int, error)
}

type ConnectionFactory func(string, string, time.Duration) (Conn, error)
type HS1xxPlug struct {
	sync.Mutex
	addr      string
	state     bool
	cnFactory ConnectionFactory
}

func NewHS1xxPlug(addr string) *HS1xxPlug {
	return &HS1xxPlug{
		addr: addr,
		cnFactory: func(proto, addr string, t time.Duration) (Conn, error) {
			return net.DialTimeout(proto, addr, t)
		},
	}
}

func (p *HS1xxPlug) On() error {
	if err := p.command(_cmdOn); err != nil {
		return err
	}
	p.state = true
	return nil
}

func (p *HS1xxPlug) Off() error {
	if err := p.command(_cmdOff); err != nil {
		return err
	}
	p.state = false
	return nil
}

func (p *HS1xxPlug) Info() error {
	return p.command(_cmdInfo)
}

func (p *HS1xxPlug) command(cmd []byte) error {
	p.Lock()
	defer p.Unlock()
	conn, err := p.cnFactory("tcp", p.addr, _timeOut)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(_timeOut)); err != nil {
		return err
	}
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(cmd)))
	bs := append(header, autokeyeEncrypt(cmd)...)
	_, err = conn.Write(bs)
	return err
}

func autokeyeEncrypt(cmd []byte) []byte {
	n := len(cmd)
	key := byte(0xAB)
	payload := make([]byte, n)
	for i := 0; i < n; i++ {
		payload[i] = cmd[i] ^ key
		key = payload[i]
	}
	return payload
}
