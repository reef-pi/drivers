package tplink

import (
	"testing"

	"github.com/reef-pi/hal"
)

func TestHS300Strip(t *testing.T) {
	meta := hal.Metadata{
		Name:        "tplink-hs300",
		Description: "tplink hs300 series smart power strip driver with current monitoring",
		Capabilities: []hal.Capability{
			hal.DigitalOutput, hal.AnalogInput,
		},
	}

	d := NewHS300Strip("127.0.0.1:9999", meta)
	nop := NewNop()
	d.SetFactory(nop.Factory)
	if d.Metadata().Name == "" {
		t.Error("HAL metadata should not have empty name")
	}

}
