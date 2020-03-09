package shelly

import (
	"fmt"
	"github.com/reef-pi/hal"
	"net/http"
)

const (
	_shelly25 = "shelly25"
)

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
