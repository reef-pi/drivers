package shelly

import (
	"errors"
	"fmt"
	"github.com/reef-pi/hal"
	"net/http"
)

type Shelly1 struct {
	meta hal.Metadata
	pins []*Relay
}

func NewShelly1(a string, devMode bool) (hal.DigitalOutputDriver, error) {
	addr := "http://" + a
	var getter HTTPGetter
	if devMode {
		getter = func(_ string) (*http.Response, error) {
			resp := new(http.Response)
			resp.StatusCode = http.StatusOK
			return resp, nil
		}
	}

	return &Shelly1{
		meta: hal.Metadata{
			Name:         "Shelly1",
			Description:  "Shelly 1, single relay wifi driver",
			Capabilities: []hal.Capability{hal.DigitalOutput},
		},
		pins: []*Relay{
			NewRelay("Shelly One Relay 0", addr, 0, getter),
		},
	}, nil
}

func (s *Shelly1) Metadata() hal.Metadata {
	return s.meta
}
func (s *Shelly1) Close() error {
	return nil
}

func (s *Shelly1) Pins(cap hal.Capability) ([]hal.Pin, error) {
	switch cap {
	case hal.DigitalOutput:
		return []hal.Pin{s.pins[0]}, nil
	default:
		return nil, fmt.Errorf("capability not supported")
	}
}
func (s *Shelly1) DigitalOutputPins() []hal.DigitalOutputPin {
	return []hal.DigitalOutputPin{s.pins[0]}
}
func (s *Shelly1) DigitalOutputPin(pin int) (hal.DigitalOutputPin, error) {
	if pin == 0 {
		return s.pins[pin], nil
	}
	return nil, fmt.Errorf("unknown pin:%d", pin)
}

type Shelly1Factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
	devMode    bool
}

func Shelly1Adapter(devMode bool) hal.DriverFactory {
	return &Shelly1Factory{
		meta: hal.Metadata{
			Name:         "Shelly1",
			Description:  "Shelly 1, single relay wifi driver",
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

func (f *Shelly1Factory) Metadata() hal.Metadata {
	return f.meta
}
func (f *Shelly1Factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *Shelly1Factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

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

func (f *Shelly1Factory) NewDriver(params map[string]interface{}, _ interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(params); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}
	addr := params[_addr].(string)
	return NewShelly1(addr, f.devMode)
}
