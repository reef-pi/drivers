package tplink

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

type HS103Plug struct {
	addr      string
	state     bool
	cnFactory ConnectionFactory
	meta      hal.Metadata
}

func NewHS103Plug(addr string) *HS103Plug {
	return &HS103Plug{
		addr: addr,
		meta: hal.Metadata{
			Name:        "tplink-hs103",
			Description: "tplink hs103 series smart plug driver",
			Capabilities: []hal.Capability{
				hal.Output,
			},
		},
		cnFactory: func(proto, addr string, t time.Duration) (Conn, error) {
			return net.DialTimeout(proto, addr, t)
		},
	}
}

func (p *HS103Plug) On() error {
	cmd := new(CmdRelayState)
	cmd.System.RelayState.State = 1
	if _, err := command(p.cnFactory, p.addr, cmd); err != nil {
		return err
	}
	p.state = true
	return nil
}

func (p *HS103Plug) Off() error {
	cmd := new(CmdRelayState)
	cmd.System.RelayState.State = 0
	if _, err := command(p.cnFactory, p.addr, cmd); err != nil {
		return err
	}
	p.state = false
	return nil
}

func (p *HS103Plug) Info() (*Sysinfo, error) {
	buf, err := command(p.cnFactory, p.addr, new(Plug))
	if err != nil {
		return nil, err
	}
	var d Plug
	if err := json.Unmarshal(buf, &d); err != nil {
		return nil, err
	}
	return &d.System.Sysinfo, nil
}

func HS103HALAdapter(c []byte, _ i2c.Bus) (hal.Driver, error) {
	var conf Config
	if err := json.Unmarshal(c, &conf); err != nil {
		return nil, err
	}
	return NewHS103Plug(conf.Address), nil
}

func (p *HS103Plug) Metadata() hal.Metadata {
	return p.meta
}

func (p *HS103Plug) Name() string {
	return p.meta.Name
}

func (p *HS103Plug) OutputPins() []hal.OutputPin {
	return []hal.OutputPin{p}
}

func (p *HS103Plug) OutputPin(i int) (hal.OutputPin, error) {
	if i != 0 {
		return nil, fmt.Errorf("invalid pin: %d", i)
	}
	return p, nil
}

func (p *HS103Plug) Write(state bool) error {
	if state {
		return p.On()
	}
	return p.Off()
}

func (p *HS103Plug) LastState() bool {
	return p.state
}

func (p *HS103Plug) Close() error {
	return nil
}
