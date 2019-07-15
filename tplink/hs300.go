package tplink

import (
	"encoding/json"
	"fmt"
	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

type (
	HS300EmeterCmd struct {
		Emeter struct {
			Realtime HS300Realtime `json:"get_realtime"`
		} `json:"emeter"`
		Context struct {
			Children []string `json:"child_ids,omitempty"`
		} `json:"context,omitempty"`
	}
	HS300Realtime struct {
		Current  float64 `json:"current_ma,omitempty"`
		Voltage  float64 `json:"voltage_mv,omitempty"`
		Power    float64 `json:"power_mw,omitempty"`
		Total    float64 `json:"total_wh,omitempty"`
		ErrrCode int     `json:"err_code,omitempty"`
	}
	Outlet struct {
		name      string
		id        string
		addr      string
		state     bool
		cnFactory ConnectionFactory
	}
)

func (o *Outlet) Name() string {
	return o.name
}

func (o *Outlet) Write(state bool) error {
	if state {
		return o.On()
	}
	return o.Off()
}

func (o *Outlet) RTEmeter() (*HS300Realtime, error) {
	var cmd HS300EmeterCmd
	cmd.Context.Children = []string{o.id}
	d, err := command(o.cnFactory, o.addr, &cmd)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(d, &cmd); err != nil {
		return nil, err
	}
	return &cmd.Emeter.Realtime, nil
}

func (o *Outlet) LastState() bool {
	return o.state
}

func (o *Outlet) On() error {
	cmd := new(CmdRelayState)
	cmd.System.RelayState.State = 1
	cmd.Context.Children = []string{o.id}
	if _, err := command(o.cnFactory, o.addr, cmd); err != nil {
		return err
	}
	o.state = true
	return nil
}
func (o *Outlet) Off() error {
	cmd := new(CmdRelayState)
	cmd.System.RelayState.State = 0
	cmd.Context.Children = []string{o.id}
	if _, err := command(o.cnFactory, o.addr, cmd); err != nil {
		return err
	}
	o.state = true
	return nil
}

func (o *Outlet) Close() error {
	return nil
}

type HS300Strip struct {
	addr      string
	cnFactory ConnectionFactory
	meta      hal.Metadata
	children  []*Outlet
}

func NewHS300Strip(addr string) *HS300Strip {
	return &HS300Strip{
		addr: addr,
		meta: hal.Metadata{
			Name:        "tplink-hs300",
			Description: "tplink hs300 series smart power strip driver with current monitoring",
			Capabilities: []hal.Capability{
				hal.Output,
			},
		},
		cnFactory: TCPConnFactory,
		children:  make([]*Outlet, 6),
	}
}

func HS300HALAdapter(c []byte, _ i2c.Bus) (hal.Driver, error) {
	var conf Config
	if err := json.Unmarshal(c, &conf); err != nil {
		return nil, err
	}
	return NewHS300Strip(conf.Address), nil
}

func (s *HS300Strip) Metadata() hal.Metadata {
	return s.meta
}

func (s *HS300Strip) Name() string {
	return s.meta.Name
}

func (s *HS300Strip) OutputPins() []hal.OutputPin {
	var pins []hal.OutputPin
	for _, o := range s.children {
		pins = append(pins, o)
	}
	return pins
}

func (s *HS300Strip) OutputPin(i int) (hal.OutputPin, error) {
	if i < 0 || i > 5 {
		return nil, fmt.Errorf("invalid pin: %d", i)
	}
	if s.children[i] == nil {
		if err := s.FetchSysInfo(); err != nil {
			return nil, err
		}
	}
	return s.children[i], nil
}

func (s *HS300Strip) Close() error {
	return nil
}
func (s *HS300Strip) FetchSysInfo() error {
	buf, err := command(s.cnFactory, s.addr, new(Plug))
	if err != nil {
		return err
	}
	var d Plug
	if err := json.Unmarshal(buf, &d); err != nil {
		return err
	}
	var children []*Outlet
	for _, ch := range d.System.Sysinfo.Children {
		o := &Outlet{
			name:      ch.Alias,
			id:        ch.ID,
			addr:      s.addr,
			cnFactory: s.cnFactory,
		}
		children = append(children, o)
	}
	s.children = children
	return nil
}

func (s *HS300Strip) Children() []*Outlet {
	return s.children
}
