package pico_board

import (
	"encoding/json"
	"fmt"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

var driverMeta = hal.Metadata{
	Name:         "pico-board",
	Description:  "Isolated ATSAMD10 pH driver on the blueAcro Pico board",
	Capabilities: []hal.Capability{hal.AnalogInput},
}

type Config struct {
	Address byte `json:"address"`
}

type driver struct {
	channels []hal.AnalogInputPin
	meta     hal.Metadata
}

func HalAdapter(c []byte, bus i2c.Bus) (hal.Driver, error) {
	return NewDriver(c, bus)
}

func NewDriver(c []byte, bus i2c.Bus) (hal.AnalogInputDriver, error) {
	var config Config
	if err := json.Unmarshal(c, &config); err != nil {
		return nil, err
	}

	ch, err := NewChannel(bus, config.Address)
	if err != nil {
		return nil, err
	}
	return &driver{
		channels: []hal.AnalogInputPin{ch},
		meta:     driverMeta,
	}, nil
}

func (d *driver) Metadata() hal.Metadata {
	return d.meta
}

func (d *driver) AnalogInputPins() []hal.AnalogInputPin {
	return d.channels
}

func (d *driver) Pins(cap hal.Capability) ([]hal.Pin, error) {
	switch cap {
	case hal.PWM, hal.DigitalOutput:
		var pins []hal.Pin
		for _, pin := range d.channels {
			pins = append(pins, pin)
		}
		return pins, nil
	default:
		return nil, fmt.Errorf("unsupported capability:%s", cap.String())
	}
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
