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
	channels []hal.ADCChannel
	meta     hal.Metadata
}

func HalAdapter(c []byte, bus i2c.Bus) (hal.Driver, error) {
	return NewDriver(c, bus)
}
func NewDriver(c []byte, bus i2c.Bus) (hal.ADCDriver, error) {
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
		channels: []hal.ADCChannel{ch},
		meta: hal.Metadata{
			Name:         "ph-board",
			Description:  "An ADS115 based analog to digital converted with onboard female BNC connector",
			Capabilities: []hal.Capability{hal.PH},
		},
	}, nil
}
func (d *driver) Metadata() hal.Metadata {
	return d.meta
}

func (d *driver) ADCChannels() []hal.ADCChannel {
	return d.channels
}

func (d *driver) ADCChannel(n int) (hal.ADCChannel, error) {
	if n != 0 {
		return nil, fmt.Errorf("ph board does not have channel %d", n)
	}
	return d.channels[0], nil
}

func (d *driver) Close() error {
	return nil
}
