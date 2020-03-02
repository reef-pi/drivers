package tplink

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/reef-pi/hal"
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
		meta     hal.Metadata
		children []*Outlet
		command  *cmd
	}
)

func NewHS300Strip(addr string, meta hal.Metadata) *HS300Strip {
	return &HS300Strip{
		meta: meta,
		command: &cmd{
			cf:   TCPConnFactory,
			addr: addr,
		},
		children: make([]*Outlet, 6),
	}
}

func (s *HS300Strip) Metadata() hal.Metadata {
	return s.meta
}

func (s *HS300Strip) SetFactory(cf ConnectionFactory) {
	s.command.cf = cf
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
	buf, err := s.command.Execute(new(Plug), true)
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
			name:    ch.Alias,
			id:      ch.ID,
			command: s.command,
			number:  i,
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

type hs300Factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var factory300 *hs300Factory
var hs300once sync.Once

// HS300Factory returns a singleton HS300 Driver factory
func HS300Factory() hal.DriverFactory {

	hs300once.Do(func() {
		factory300 = &hs300Factory{
			meta: hal.Metadata{
				Name:        "tplink-hs300",
				Description: "tplink hs300 series smart power strip driver with current monitoring",
				Capabilities: []hal.Capability{
					hal.DigitalOutput, hal.AnalogInput,
				},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    addressParam,
					Type:    hal.String,
					Order:   0,
					Default: "192.168.1.11:9999",
				},
			},
		}
	})

	return factory300
}

func (f *hs300Factory) Metadata() hal.Metadata {
	return f.meta
}

func (f *hs300Factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *hs300Factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

	var failures = make(map[string][]string)

	if v, ok := parameters[addressParam]; ok {
		_, ok := v.(string)
		if !ok {
			failure := fmt.Sprint(addressParam, " is not a string. ", v, " was received.")
			failures[addressParam] = append(failures[addressParam], failure)
		}
	} else {
		failure := fmt.Sprint(addressParam, " is a required parameter, but was not received.")
		failures[addressParam] = append(failures[addressParam], failure)
	}

	return len(failures) == 0, failures
}

func (f *hs300Factory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	addr := parameters[addressParam].(string)

	s := NewHS300Strip(addr, f.meta)
	return s, s.FetchSysInfo()
}
