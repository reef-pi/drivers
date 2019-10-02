package ph_board

import (
	"encoding/json"
	"fmt"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

const driverName = "ph-board"

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
	if err := bus.WriteBytes(config.Address, []byte{0x06}); err != nil {
		return nil, err
	}
	if err := bus.WriteBytes(config.Address, []byte{0x40, 0x06}); err != nil {
		return nil, err
	}
	if err := bus.WriteBytes(config.Address, []byte{0x08}); err != nil {
		return nil, err
	}

	ch, err := NewChannel(bus, config.Address)
	if err != nil {
		return nil, err
	}
	return &driver{
		channels: []hal.AnalogInputPin{ch},
		meta: hal.Metadata{
			Name:         "ph-board",
			Description:  "An ADS115 based analog to digital converted with onboard female BNC connector",
			Capabilities: []hal.Capability{hal.AnalogInput},
		},
	}, nil
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
