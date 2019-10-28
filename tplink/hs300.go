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
	HS300Strip struct {
		addr      string
		cnFactory ConnectionFactory
		meta      hal.Metadata
		children  []*Outlet
	}
)

func NewHS300Strip(addr string) *HS300Strip {
	return &HS300Strip{
		addr: addr,
		meta: hal.Metadata{
			Name:        "tplink-hs300",
			Description: "tplink hs300 series smart power strip driver with current monitoring",
			Capabilities: []hal.Capability{
				hal.DigitalOutput, hal.AnalogInput,
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
	s := NewHS300Strip(conf.Address)
	return s, s.FetchSysInfo()
}

func (s *HS300Strip) Metadata() hal.Metadata {
	return s.meta
}

func (s *HS300Strip) Name() string {
	return s.meta.Name
}

func (s *HS300Strip) DigitalOutputPins() []hal.DigitalOutputPin {
	var pins []hal.DigitalOutputPin
	for _, o := range s.children {
		pins = append(pins, o)
	}
	return pins
}

func (s *HS300Strip) DigitalOutputPin(i int) (hal.DigitalOutputPin, error) {
	if i < 0 || i > 5 {
		return nil, fmt.Errorf("invalid pin: %d", i)
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
		fmt.Println(string(buf))
		return err
	}
	var children []*Outlet
	for i, ch := range d.System.Sysinfo.Children {
		o := &Outlet{
			name:      ch.Alias,
			id:        ch.ID,
			addr:      s.addr,
			cnFactory: s.cnFactory,
			number:    i,
		}
		children = append(children, o)
	}
	s.children = children
	return nil
}

func (s *HS300Strip) Children() []*Outlet {
	return s.children
}

func (p *HS300Strip) AnalogInputPins() []hal.AnalogInputPin {
	var channels []hal.AnalogInputPin
	for _, o := range p.children {
		channels = append(channels, o)
	}
	return channels
}

func (p *HS300Strip) AnalogInputPin(i int) (hal.AnalogInputPin, error) {
	if i < 0 || i > 5 {
		return nil, fmt.Errorf("invalid channel number: %d", i)
	}
	return p.children[i], nil
}

func (p *HS300Strip) Pins(cap hal.Capability) ([]hal.Pin, error) {
	switch cap {
	case hal.DigitalOutput, hal.AnalogInput:
		var channels []hal.Pin
		for _, o := range p.children {
			channels = append(channels, o)
		}
		return channels, nil
	default:
		return nil, fmt.Errorf("unsupported capability:%s", cap.String())
	}
}
