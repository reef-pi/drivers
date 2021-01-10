package ads1x15

import (
	"fmt"

	"github.com/reef-pi/hal"
)

type driver struct {
	channels []hal.AnalogInputPin
	meta     hal.Metadata
}

func (d *driver) Metadata() hal.Metadata {
	return d.meta
}

func (d *driver) Pins(cap hal.Capability) ([]hal.Pin, error) {
	if cap == hal.AnalogInput {
		pins := make([]hal.Pin, len(d.channels))
		for i, ch := range d.channels {
			pins[i] = ch
		}
		return pins, nil
	}

	return nil, fmt.Errorf("unsupported capability: %s", cap.String())
}

func (d *driver) AnalogInputPins() []hal.AnalogInputPin {
	return d.channels
}

func (d *driver) AnalogInputPin(n int) (hal.AnalogInputPin, error) {
	if n < 0 || n > 3 {
		return nil, fmt.Errorf("ADS1X15 does not have channel %d", n)
	}
	return d.channels[n], nil
}

func (d *driver) Close() error {
	return nil
}
