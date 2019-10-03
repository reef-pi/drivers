package pico_board

import (
	"testing"

	"github.com/reef-pi/hal"
	"github.com/reef-pi/rpi/i2c"
)

func TestPhBoardDriver(t *testing.T) {
	bus := i2c.MockBus()
	bus.Bytes = make([]byte, 2)
	_, err := HalAdapter([]byte(""), bus)
	if err == nil {
		t.Error("Adapter creation should fail when json config is invalid")
	}
	configJSON := `
	{
		"address":72
	}
	`
	d, err := NewDriver([]byte(configJSON), bus)
	if err != nil {
		t.Error(err)
	}
	if d.Metadata().Name != "pico-board" {
		t.Error("Unexpected name")
	}
	if !d.Metadata().HasCapability(hal.AnalogInput) {
		t.Error("analog input cpability should exist")
	}
	if d.Metadata().HasCapability(hal.DigitalInput) {
		t.Error("Digital input Capability should not exist")
	}

	if len(d.AnalogInputPins()) != 1 {
		t.Error("Expected only one channel")
	}
	if _, err := d.AnalogInputPin(1); err == nil {
		t.Error("Expected error for invalid channel name")
	}

	ch, err := d.AnalogInputPin(0)
	if err != nil {
		t.Error(err)
	}
	if ch.Name() != "0" {
		t.Error("Unexpected channel name")
	}
	v, err := ch.Read()
	if err != nil {
		t.Error(err)
	}
	if v != 0 {
		t.Error("Unexepected value")
	}
	if err := d.Close(); err != nil {
		t.Error(err)
	}
}
