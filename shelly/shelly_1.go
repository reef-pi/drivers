package shelly

import (
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
			return new(http.Response), nil
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
		return nil, fmt.Errorf("unspoorted capability")
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
