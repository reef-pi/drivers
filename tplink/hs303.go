package tplink

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/reef-pi/hal"
)

type HS303Strip struct {
	meta     hal.Metadata
	children []*Outlet
	command  *cmd
}

func NewHS303Strip(addr string, meta hal.Metadata) *HS303Strip {
	return &HS303Strip{
		meta: meta,
		command: &cmd{
			cf:   TCPConnFactory,
			addr: addr,
		},
		children: make([]*Outlet, 3),
	}
}

func (s *HS303Strip) Metadata() hal.Metadata {
	return s.meta
}

func (s *HS303Strip) SetFactory(cf ConnectionFactory) {
	s.command.cf = cf
}
func (s *HS303Strip) Name() string {
	return s.meta.Name
}

func (s *HS303Strip) DigitalOutputPins() []hal.DigitalOutputPin {
	var pins []hal.DigitalOutputPin
	for _, o := range s.children {
		pins = append(pins, o)
	}
	return pins
}

func (s *HS303Strip) DigitalOutputPin(i int) (hal.DigitalOutputPin, error) {
	if i < 0 || i > 2 {
		return nil, fmt.Errorf("invalid pin: %d", i)
	}
	return s.children[i], nil
}

func (s *HS303Strip) Close() error {
	return nil
}
func (s *HS303Strip) FetchSysInfo() error {
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

func (s *HS303Strip) Children() []*Outlet {
	return s.children
}

func (p *HS303Strip) Pins(cap hal.Capability) ([]hal.Pin, error) {
	switch cap {
	case hal.DigitalOutput:
		var channels []hal.Pin
		for _, o := range p.children {
			channels = append(channels, o)
		}
		return channels, nil
	default:
		return nil, fmt.Errorf("unsupported capability:%s", cap.String())
	}
}

type hs303Factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

var factory303 *hs303Factory
var hs303once sync.Once

// HS303Factory returns a singleton HS300 Driver factory
func HS303Factory() hal.DriverFactory {

	hs303once.Do(func() {
		factory303 = &hs303Factory{
			meta: hal.Metadata{
				Name:        "tplink-hs303",
				Description: "tplink hs303 series smart power strip driver",
				Capabilities: []hal.Capability{
					hal.DigitalOutput,
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

	return factory303
}

func (f *hs303Factory) Metadata() hal.Metadata {
	return f.meta
}

func (f *hs303Factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *hs303Factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

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

func (f *hs303Factory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	addr := parameters[addressParam].(string)

	s := NewHS303Strip(addr, f.meta)
	return s, s.FetchSysInfo()
}
