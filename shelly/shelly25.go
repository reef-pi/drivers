package shelly

import (
	"errors"
	"fmt"
	"github.com/reef-pi/hal"
	"net/http"
)

const (
	_shelly25 = "shelly25"
	_addr     = "Address"
)

type Shelly25Factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
	devMode    bool
}

func Shelly25Adapter(devMode bool) hal.DriverFactory {
	return &Shelly25Factory{
		meta: hal.Metadata{
			Name:         "Shelly2,5",
			Description:  "Shelly 2.5 , dual relay wifi driver",
			Capabilities: []hal.Capability{hal.DigitalOutput},
		},
		parameters: []hal.ConfigParameter{
			{
				Name:    _addr,
				Type:    hal.String,
				Order:   0,
				Default: "192.168.1.33",
			},
		},
		devMode: devMode,
	}
}

func (f *Shelly25Factory) Metadata() hal.Metadata {
	return f.meta
}
func (f *Shelly25Factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *Shelly25Factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

	var failures = make(map[string][]string)

	if v, ok := parameters[_addr]; ok {
		_, ok := v.(string)
		if !ok {
			failure := fmt.Sprint(_addr, " is not a string. ", v, " was received.")
			failures[_addr] = append(failures[_addr], failure)
		}
	} else {
		failure := fmt.Sprint(_addr, " is a required parameter, but was not received.")
		failures[_addr] = append(failures[_addr], failure)
	}

	return len(failures) == 0, failures
}

func (f *Shelly25Factory) NewDriver(params map[string]interface{}, _ interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(params); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}
	addr := params[_addr].(string)
	return NewShelly25(addr, f.devMode)
}

type Shelly25 struct {
	meta hal.Metadata
	pins []*Relay
}

func NewShelly25(a string, devMode bool) (hal.DigitalOutputDriver, error) {
	addr := "http://" + a
	var getter HTTPGetter
	if devMode {
		getter = func(_ string) (*http.Response, error) {
			return new(http.Response), nil
		}
	}

	return &Shelly25{
		meta: hal.Metadata{
			Name:         "Shelly2,5",
			Description:  "Shelly 2.5 , dual relay wifi driver",
			Capabilities: []hal.Capability{hal.DigitalOutput},
		},
		pins: []*Relay{
			NewRelay("Shelly 2.5 Relay 0", addr, 0, getter),
			NewRelay("Shelly 2.5 Relay 1", addr, 1, getter),
		},
	}, nil
}

func (s *Shelly25) Metadata() hal.Metadata {
	return s.meta
}
func (s *Shelly25) Close() error {
	return nil
}

func (s *Shelly25) Pins(cap hal.Capability) ([]hal.Pin, error) {
	switch cap {
	case hal.DigitalOutput:
		return []hal.Pin{s.pins[0], s.pins[1]}, nil
	default:
		return nil, fmt.Errorf("unspoorted capability")
	}
}
func (s *Shelly25) DigitalOutputPins() []hal.DigitalOutputPin {
	return []hal.DigitalOutputPin{s.pins[0], s.pins[1]}
}
func (s *Shelly25) DigitalOutputPin(pin int) (hal.DigitalOutputPin, error) {
	if pin == 0 || pin == 1 {
		return s.pins[pin], nil
	}
	return nil, fmt.Errorf("unknown pin:%d", pin)
}
