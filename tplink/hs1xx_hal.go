package tplink

import (
	"encoding/json"
	"fmt"
	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

type HS1xxPlugConfig struct {
	Address string `json:"address"`
}

var meta = hal.Metadata{
	Name:        "tplink-hs1xx",
	Description: "Supports tplink hs1xx series smart plugs",
	Capabilities: []hal.Capability{
		hal.Output,
	},
}

func HALAdapter(c []byte, _ i2c.Bus) (hal.Driver, error) {
	var conf HS1xxPlugConfig
	if err := json.Unmarshal(c, &conf); err != nil {
		return nil, err
	}
	return NewHS1xxPlug(conf.Address), nil
}

func (p *HS1xxPlug) Metadata() hal.Metadata {
	return meta
}
func (p *HS1xxPlug) Name() string {
	return meta.Name
}

func (p *HS1xxPlug) OutputPins() []hal.OutputPin {
	return []hal.OutputPin{p}
}
func (p *HS1xxPlug) OutputPin(i int) (hal.OutputPin, error) {
	if i != 0 {
		return nil, fmt.Errorf("invalid pin: %d", i)
	}
	return p, nil
}

func (p *HS1xxPlug) Write(state bool) error {
	if state {
		return p.On()
	}
	return p.Off()
}

func (p *HS1xxPlug) LastState() bool {
	return p.state
}
func (p *HS1xxPlug) Close() error {
	return nil
}
