package dli

import (
	"fmt"
	"github.com/reef-pi/hal"
)

type Driver struct {
	meta   hal.Metadata
	relays []*Relay
}

func NewDriver(a, u, p string) *Driver {
	conf := Config{
		addr:     a,
		username: u,
		password: p,
	}
	return &Driver{
		meta: hal.Metadata{
			Name:         "DLI-Webpowerswitch-Pro",
			Description:  "DLI Web power switch pro",
			Capabilities: []hal.Capability{hal.DigitalOutput},
		},
		relays: []*Relay{
			&Relay{0, conf, false},
			&Relay{1, conf, false},
			&Relay{2, conf, false},
			&Relay{3, conf, false},
			&Relay{4, conf, false},
			&Relay{5, conf, false},
			&Relay{6, conf, false},
			&Relay{7, conf, false},
		},
	}
}
func (d *Driver) Metadata() hal.Metadata {
	return d.meta
}
func (d *Driver) Close() error {
	return nil
}

func (d *Driver) Pins(c hal.Capability) ([]hal.Pin, error) {
	switch c {
	case hal.DigitalOutput:
		return []hal.Pin{
			d.relays[0],
			d.relays[1],
			d.relays[2],
			d.relays[3],
			d.relays[4],
			d.relays[5],
			d.relays[6],
			d.relays[7],
		}, nil
	default:
		return nil, fmt.Errorf("capability not supported")
	}
}
func (d *Driver) DigitalOutputPins() []hal.DigitalOutputPin {
	return []hal.DigitalOutputPin{
		d.relays[0],
		d.relays[1],
		d.relays[2],
		d.relays[3],
		d.relays[4],
		d.relays[5],
		d.relays[6],
		d.relays[7],
	}
}
func (d *Driver) DigitalOutputPin(pin int) (hal.DigitalOutputPin, error) {
	if pin >= 0 && pin < 8 {
		return d.relays[pin], nil
	}
	return nil, fmt.Errorf("unknown pin:%d", pin)
}
