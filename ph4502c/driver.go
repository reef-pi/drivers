package ph4502c

import (
	"fmt"

	"github.com/reef-pi/hal"
)

const driverName = "ph4502c"

type driver struct {
	channels []hal.AnalogInputPin
	meta     hal.Metadata
}

func (d *driver) Metadata() hal.Metadata {
	return d.meta
}
func (d *driver) Pins(cap hal.Capability) ([]hal.Pin, error) {
	if cap == hal.AnalogInput {
		return []hal.Pin{d.channels[0]}, nil
	}
	return nil, fmt.Errorf("unsupported capability: %s", cap.String())
}

func (d *driver) AnalogInputPins() []hal.AnalogInputPin {
	return d.channels
}

func (d *driver) AnalogInputPin(n int) (hal.AnalogInputPin, error) {
	if n != 0 {
		return nil, fmt.Errorf("ph board does not have channel %d", n)
	}
	return d.channels[0], nil
}

func (d *driver) Close() error {
	return nil
}
