package pico_board

import (
	"encoding/json"
	"fmt"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

var driverMeta =  hal.Metadata{
	Name:         "pico-board",
	Description:  "Isolated ATSAMD10 pH driver on the blueAcro Pico board",
	Capabilities: []hal.Capability{hal.PH},
}

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

	ch, err := NewChannel(bus, config.Address)
	if err != nil {
		return nil, err
	}
	return &driver{
		channels: []hal.ADCChannel{ch},
		meta: driverMeta,
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
